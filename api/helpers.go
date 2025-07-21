package api

import (
	"cold_emailer/apollo"
	"cold_emailer/constants"
	"cold_emailer/dbmodels/companies"
	"cold_emailer/dbmodels/contacts"
	"cold_emailer/dbmodels/email_logs"
	"cold_emailer/dbmodels/gmailtokens"
	"cold_emailer/dbmodels/profileinfo"
	"cold_emailer/gmail"
	"cold_emailer/openai"
	"context"
	"math"
	"time"

	"cold_emailer/utils"

	"github.com/google/uuid"
)

func EnrichDBWithCompaniesAndContacts(ctx context.Context, countNewCompanies, maxContactsPerCompany int) (err error) {
	client := apollo.NewApolloClient()

	pageNum := 1
	maxPages := math.MaxInt

	for countNewCompanies > 0 && pageNum <= maxPages {
		resp, err := client.SearchCompaniesMixed(apollo.MixedCompanySearchRequest{
			OrganizationNumEmployeesRanges: []string{"1,10", "11,50", "51,200", "201,500"},
			OrganizationLocations:          constants.TargetCountries,
			QOrganizationKeywordTags:       []string{"technology", "software", "IT", "artificial intelligence"},
			Page:                           pageNum,
			PerPage:                        10,
		})
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return err
		}
		if pageNum == 1 {
			maxPages = resp.Pagination.TotalPages
			utils.Logger.Info().Int("maxPages", maxPages).Msg("maxPages for company search:")
		}

		if len(resp.Accounts) == 0 {
			utils.Logger.Info().Msg("No more companies to add")
			return nil
		}

		for _, fetchedCompany := range resp.Accounts {
			utils.Logger.Info().Str("company_name", fetchedCompany.Name).Str("organization_id", fetchedCompany.OrganizationID).Msg("Checking company:")
			if !fetchedCompany.CheckBasicFields() {
				utils.Logger.Info().Str("company_name", fetchedCompany.Name).Str("organization_id", fetchedCompany.OrganizationID).Msg("Skipping this company, failing basic fields check...")
				continue
			}

			exists, err := companies.ExistsByApolloID(ctx, fetchedCompany.OrganizationID, companies.StatusActive)
			if err != nil {
				utils.Logger.Error().Err(err).Msg("err:")
				return err
			}

			if exists {
				utils.Logger.Info().Str("company_name", fetchedCompany.Name).Str("organization_id", fetchedCompany.OrganizationID).Msg("Skipping company because it already exists")
				continue
			}

			apolloCompanyDetailsResp, err := client.GetOrganizationDetails(fetchedCompany.OrganizationID)
			if err != nil {
				utils.Logger.Error().Err(err).Msg("error:")
				return err
			}

			dbCompanyID, err := companies.InsertIfNotExists(ctx, companies.CompanyForSet{
				ID:             uuid.New().String(),
				ApolloID:       fetchedCompany.OrganizationID,
				Status:         companies.StatusActive,
				Name:           fetchedCompany.Name,
				Website:        fetchedCompany.WebsiteURL,
				Industry:       apolloCompanyDetailsResp.Organization.Industry,
				SubIndustry:    "",
				TechDetails:    apolloCompanyDetailsResp.Organization.StringifyTechnologyArray(),
				CompanyDetails: apolloCompanyDetailsResp.Organization.ShortDescription,
				Metadata:       "{}",
			})

			if err != nil {
				utils.Logger.Error().Err(err).Msg("err:")
				return err
			}

			err = EnrichContactsForOrganizationID(ctx, fetchedCompany.OrganizationID, dbCompanyID, maxContactsPerCompany)
			if err != nil {
				if err == constants.ErrorNoPeopleFound {
					utils.Logger.Info().Str("company_id", dbCompanyID).Msg("no contacts found for the company, deleting it")
					err = companies.DeleteByID(ctx, dbCompanyID)
					if err != nil {
						utils.Logger.Error().Err(err).Msg("err:")
						return err
					}
					continue
				}
				utils.Logger.Error().Err(err).Msg("err:")
				return err
			}

			countNewCompanies--
			if countNewCompanies <= 0 {
				break
			}
		}
		pageNum++
	}

	return nil
}

func EnrichContactsForOrganizationID(ctx context.Context, apolloOrganizationID, dbCompanyID string, maxContactsForOrg int) (err error) {
	client := apollo.NewApolloClient()
	pageNum := 1
	maxPages := math.MaxInt

	for maxContactsForOrg > 0 && pageNum <= maxPages {
		resp, err := client.SearchContacts(apollo.MixedPeopleSearchRequest{
			PersonTitles:    constants.SuitableRoles,
			OrganizationIDs: []string{apolloOrganizationID},
			Page:            pageNum,
			PerPage:         10,
		})
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return err
		}

		if pageNum == 1 {
			if len(resp.People) == 0 {
				err = constants.ErrorNoPeopleFound
				utils.Logger.Warn().Err(err).Msg("warn err:")
				return err
			}
			maxPages = resp.Pagination.TotalPages
			utils.Logger.Info().Int("maxPages", maxPages).Msg("maxPages for this company's person search:")
		}

		if len(resp.People) == 0 {
			utils.Logger.Info().Msg("No more persons to add")
			return nil
		}

		for _, fetchedPerson := range resp.People {
			utils.Logger.Info().Str("person_name", fetchedPerson.Name).Msg("checking person:")
			if !fetchedPerson.CheckBasicFields() {
				utils.Logger.Info().Str("person_name", fetchedPerson.Name).Msg("Skipping this person, failing basic fields check...")
				continue
			}

			if fetchedPerson.Email == constants.ApolloLockedEmail {
				utils.Logger.Info().Str("person_id", fetchedPerson.ID).Msg("enriching email")
				enrichContactResp, err := client.EnrichContactByApolloID(apollo.PersonMatchRequest{
					ID:                   fetchedPerson.ID,
					RevealPersonalEmails: true,
				})
				if err != nil {
					utils.Logger.Error().Err(err).Msg("err:")
					return err
				}
				utils.Logger.Info().Str("enriched_email", enrichContactResp.Person.Email).Msg("enriched email:")
				if enrichContactResp.Person.Email == "" {
					utils.Logger.Info().Str("person_id", fetchedPerson.ID).Msg("enriched email is empty, skipping this person")
					continue
				}
				fetchedPerson.Email = enrichContactResp.Person.Email
			}

			exists, err := contacts.ExistsByApolloID(ctx, fetchedPerson.ID, contacts.StatusActive)
			if err != nil {
				utils.Logger.Error().Err(err).Msg("err:")
				return err
			}

			if exists {
				utils.Logger.Info().Str("person_id", fetchedPerson.ID).Msg("contact already exists, skipping...")
				continue
			}

			_, err = contacts.InsertIfNotExists(ctx, contacts.ContactForSet{
				ID:          uuid.New().String(),
				ApolloID:    fetchedPerson.ID,
				CompanyID:   dbCompanyID,
				Status:      contacts.StatusActive,
				Name:        fetchedPerson.Name,
				EmailID:     fetchedPerson.Email,
				LinkedInURL: fetchedPerson.LinkedinURL,
				PhoneNumber: fetchedPerson.Organization.SanitizedPhone,
				Role:        fetchedPerson.Title,
				Metadata:    "{}",
			})

			if err != nil {
				utils.Logger.Error().Err(err).Msg("err:")
				return err
			}

			maxContactsForOrg--
			if maxContactsForOrg <= 0 {
				break
			}
		}
		pageNum++
	}

	return nil
}

func GenerateAndSendEmails(ctx context.Context, limit int, status string) (err error) {
	// Get Gmail token
	token, err := gmailtokens.GetLatestToken(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
		return err
	}

	// Step 2 - Fetch contacts with company info
	contactsList, err := contacts.GetContactsWithCompanyInfo(ctx, limit, status, "")
	if err != nil {
		utils.Logger.Error().Err(err).Msgf("ERROR: Failed to fetch contacts: %v", err)
		return err
	}

	utils.Logger.Info().Int("found_contacts", len(contactsList)).Msg("Found contacts to send emails to")

	// Get profile info for email generation
	profile, err := profileinfo.GetLatestActive(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msgf("ERROR: Failed to get profile info: %v", err)
		return err
	}

	from := token.EmailID
	if from == "me" {
		from = profile.Email // Use profile email as fallback
	}

	// Initialize OpenAI client
	openaiClient := openai.NewOpenAIClient()

	// Track results
	emailsGenerated := 0
	emailsSent := 0
	emailsFailed := 0

	// Process each contact
	for _, contact := range contactsList {
		utils.Logger.Info().Str("contact_name", contact.ContactName).Str("contact_email", contact.ContactEmail).Str("company_name", contact.CompanyName).Msgf("Processing contact: %s (%s) at %s", contact.ContactName, contact.ContactEmail, contact.CompanyName)

		// Step 3: Generate personalized email
		emailData := openai.EmailGenerationData{
			ContactName:       contact.ContactName,
			ContactRole:       contact.ContactRole,
			ContactLinkedIn:   contact.ContactLinkedIn,
			CompanyName:       contact.CompanyName,
			CompanyWebsite:    contact.CompanyWebsite,
			CompanyIndustry:   contact.CompanyIndustry,
			CompanyDetails:    contact.CompanyDetails,
			CompanyTechStack:  contact.CompanyTech,
			ProfileName:       profile.Name,
			ProfileExperience: profile.Experience,
			ProfileSkills:     profile.Skills,
			ProfileSummary:    profile.Summary,
		}

		subject, body, err := openaiClient.GeneratePersonalizedEmail(emailData)
		if err != nil {
			utils.Logger.Error().Err(err).Msgf("ERROR: Failed to generate email for %s: %v", contact.ContactEmail, err)
			emailsFailed++
			continue
		}

		emailsGenerated++
		utils.Logger.Info().Str("contact_email", contact.ContactEmail).Str("subject", subject).Msgf("Generated email for %s: %s", contact.ContactEmail, subject)

		// Log GENERATED stage
		metadata := map[string]interface{}{
			"contact_name": contact.ContactName,
			"contact_role": contact.ContactRole,
			"company_name": contact.CompanyName,
			"generated_at": time.Now().Format(time.RFC3339),
		}

		if err := email_logs.LogGenerated(ctx, contact.ContactID, contact.CompanyID, subject, body, metadata); err != nil {
			utils.Logger.Error().Err(err).Msgf("ERROR: Failed to log GENERATED stage for %s: %v", contact.ContactEmail, err)
		}

		// Step 4: Send email with resume
		if profile.ResumePath == "" {
			utils.Logger.Info().Str("contact_email", contact.ContactEmail).Msg("WARNING: No resume path found for profile, sending without attachment")
			err = gmail.SendSingleEmail(ctx, token.AccessToken, from, contact.ContactEmail, subject, body)
		} else {
			err = gmail.SendEmailWithAttachment(ctx, token.AccessToken, from, contact.ContactEmail, subject, body, profile.ResumePath)
		}

		if err != nil {
			utils.Logger.Error().Err(err).Msgf("ERROR: Failed to send email to %s: %v", contact.ContactEmail, err)
			emailsFailed++

			// Log ERROR stage
			errorMetadata := map[string]interface{}{
				"contact_name": contact.ContactName,
				"contact_role": contact.ContactRole,
				"company_name": contact.CompanyName,
				"error_at":     time.Now().Format(time.RFC3339),
				"error_type":   "send_failed",
			}

			if logErr := email_logs.LogError(ctx, contact.ContactID, contact.CompanyID, subject, body, err.Error(), errorMetadata); logErr != nil {
				utils.Logger.Error().Err(logErr).Msgf("ERROR: Failed to log ERROR stage for %s: %v", contact.ContactEmail, logErr)
			}
			continue
		}

		emailsSent++
		utils.Logger.Info().Str("contact_email", contact.ContactEmail).Str("subject", subject).Msgf("SUCCESS: Sent email to %s (%s)", contact.ContactEmail, subject)

		// Log SENT stage
		sentMetadata := map[string]interface{}{
			"contact_name":    contact.ContactName,
			"contact_role":    contact.ContactRole,
			"company_name":    contact.CompanyName,
			"sent_at":         time.Now().Format(time.RFC3339),
			"resume_attached": profile.ResumePath != "",
		}

		if err := email_logs.LogSent(ctx, contact.ContactID, contact.CompanyID, subject, body, sentMetadata); err != nil {
			utils.Logger.Error().Err(err).Msgf("ERROR: Failed to log SENT stage for %s: %v", contact.ContactEmail, err)
		}

		// Small delay to avoid rate limiting
		time.Sleep(2 * time.Second)
	}

	utils.Logger.Info().Int("email_sending_completed", emailsGenerated).Int("sent", emailsSent).Int("failed", emailsFailed).Msg("Email sending completed:")

	return nil
}

func GenerateAndSendFollowUpEmails(ctx context.Context, daysPastFirstEmail int, limit int, status string) (err error) {
	// Get Gmail token
	token, err := gmailtokens.GetLatestToken(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("error:")
		return err
	}

	// Step 2 - Fetch follow-up candidates
	candidates, err := contacts.GetFollowUpCandidates(ctx, daysPastFirstEmail, limit, status)
	if err != nil {
		utils.Logger.Error().Err(err).Msgf("ERROR: Failed to fetch follow-up candidates: %v", err)
		return err
	}

	utils.Logger.Info().Int("found_contacts", len(candidates)).Msg("Found contacts to send follow-up emails to")

	// Get profile info for email generation
	profile, err := profileinfo.GetLatestActive(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msgf("ERROR: Failed to get profile info: %v", err)
		return err
	}

	from := token.EmailID
	if from == "me" {
		from = profile.Email // Use profile email as fallback
	}

	// Initialize OpenAI client
	openaiClient := openai.NewOpenAIClient()

	emailsGenerated := 0
	emailsSent := 0
	emailsFailed := 0

	for _, candidate := range candidates {
		utils.Logger.Info().Str("contact_name", candidate.ContactName).Str("contact_email", candidate.ContactEmail).Str("company_name", candidate.CompanyName).Msgf("Processing follow-up for contact: %s (%s) at %s", candidate.ContactName, candidate.ContactEmail, candidate.CompanyName)

		// Step 3: Generate personalized follow-up email
		emailData := openai.EmailGenerationData{
			ContactName:       candidate.ContactName,
			ContactRole:       candidate.ContactRole,
			ContactLinkedIn:   candidate.ContactLinkedIn,
			CompanyName:       candidate.CompanyName,
			CompanyWebsite:    candidate.CompanyWebsite,
			CompanyIndustry:   candidate.CompanyIndustry,
			CompanyDetails:    candidate.CompanyDetails,
			CompanyTechStack:  candidate.CompanyTech,
			ProfileName:       profile.Name,
			ProfileExperience: profile.Experience,
			ProfileSkills:     profile.Skills,
			ProfileSummary:    profile.Summary,
		}

		// Add previous email context to the prompt
		followUpPrompt := "This is a follow-up to a previous email. Previous subject: '" + candidate.EmailSubject + "'. Previous body: '" + candidate.EmailBody + "'. Please reference the previous outreach and write a polite, concise follow-up."

		subject, body, err := openaiClient.GeneratePersonalizedEmailWithExtraPrompt(emailData, followUpPrompt)
		if err != nil {
			utils.Logger.Error().Err(err).Msgf("ERROR: Failed to generate follow-up email for %s: %v", candidate.ContactEmail, err)
			emailsFailed++
			continue
		}

		emailsGenerated++
		utils.Logger.Info().Str("contact_email", candidate.ContactEmail).Str("subject", subject).Msgf("Generated follow-up email for %s: %s", candidate.ContactEmail, subject)

		// Log GENERATED_FOLLOW_UP stage
		metadata := map[string]interface{}{
			"contact_name": candidate.ContactName,
			"contact_role": candidate.ContactRole,
			"company_name": candidate.CompanyName,
			"generated_at": time.Now().Format(time.RFC3339),
		}

		if err := email_logs.LogGeneratedFollowUp(ctx, candidate.ContactID, candidate.CompanyID, subject, body, metadata); err != nil {
			utils.Logger.Error().Err(err).Msgf("ERROR: Failed to log GENERATED_FOLLOW_UP stage for %s: %v", candidate.ContactEmail, err)
		}

		// Step 4: Send follow-up email with resume
		if profile.ResumePath == "" {
			utils.Logger.Info().Str("contact_email", candidate.ContactEmail).Msg("WARNING: No resume path found for profile, sending without attachment")
			err = gmail.SendSingleEmail(ctx, token.AccessToken, from, candidate.ContactEmail, subject, body)
		} else {
			err = gmail.SendEmailWithAttachment(ctx, token.AccessToken, from, candidate.ContactEmail, subject, body, profile.ResumePath)
		}

		if err != nil {
			utils.Logger.Error().Err(err).Msgf("ERROR: Failed to send follow-up email to %s: %v", candidate.ContactEmail, err)
			emailsFailed++

			// Log ERROR stage
			errorMetadata := map[string]interface{}{
				"contact_name": candidate.ContactName,
				"contact_role": candidate.ContactRole,
				"company_name": candidate.CompanyName,
				"error_at":     time.Now().Format(time.RFC3339),
				"error_type":   "send_failed_follow_up",
			}

			if logErr := email_logs.LogError(ctx, candidate.ContactID, candidate.CompanyID, subject, body, err.Error(), errorMetadata); logErr != nil {
				utils.Logger.Error().Err(logErr).Msgf("ERROR: Failed to log ERROR stage for %s: %v", candidate.ContactEmail, logErr)
			}
			continue
		}

		emailsSent++
		utils.Logger.Info().Str("contact_email", candidate.ContactEmail).Str("subject", subject).Msgf("SUCCESS: Sent follow-up email to %s (%s)", candidate.ContactEmail, subject)

		// Log SENT_FOLLOW_UP stage
		sentMetadata := map[string]interface{}{
			"contact_name":    candidate.ContactName,
			"contact_role":    candidate.ContactRole,
			"company_name":    candidate.CompanyName,
			"sent_at":         time.Now().Format(time.RFC3339),
			"resume_attached": profile.ResumePath != "",
		}

		if err := email_logs.LogSentFollowUp(ctx, candidate.ContactID, candidate.CompanyID, subject, body, sentMetadata); err != nil {
			utils.Logger.Error().Err(err).Msgf("ERROR: Failed to log SENT_FOLLOW_UP stage for %s: %v", candidate.ContactEmail, err)
		}

		// Small delay to avoid rate limiting
		time.Sleep(2 * time.Second)
	}

	utils.Logger.Info().Int("follow_up_email_sending_completed", emailsGenerated).Int("sent", emailsSent).Int("failed", emailsFailed).Msg("Follow-up email sending completed:")

	return nil
}

func BackFillCompanyDetailsFunc(ctx context.Context) (err error) {
	client := apollo.NewApolloClient()
	utils.Logger.Info().Msg("fetching companies from db")
	arrCompanies, err := companies.GetAll(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("err:")
		return err
	}

	utils.Logger.Info().Msg("fetching details for each company")
	for _, arrCompany := range arrCompanies {

		if arrCompany.ApolloID == "" {
			utils.Logger.Info().Str("company_name", arrCompany.Name).Msgf("for company: %s, apollo_id is empty, skipping backfill details", arrCompany.Name)
			continue
		}

		utils.Logger.Info().Str("company_name", arrCompany.Name).Str("organization_id", arrCompany.ApolloID).Msg("fetching details for company:")

		time.Sleep(1 * time.Second)
		apolloOrgDetailsResp, err := client.GetOrganizationDetails(arrCompany.ApolloID)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return err
		}

		techDetails := apolloOrgDetailsResp.Organization.StringifyTechnologyArray()
		err = companies.Update(ctx, companies.UpdateInput{
			ID:             arrCompany.ID,
			Industry:       &apolloOrgDetailsResp.Organization.Industry,
			TechDetails:    &techDetails,
			CompanyDetails: &apolloOrgDetailsResp.Organization.ShortDescription,
		})
		if err != nil {
			utils.Logger.Error().Err(err).Msg("err:")
			return err
		}
	}

	utils.Logger.Info().Msg("backfill complete")
	return nil
}

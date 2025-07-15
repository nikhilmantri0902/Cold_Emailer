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
	"log"
	"math"
	"time"

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
			log.Println("err:", err)
			return err
		}
		if pageNum == 1 {
			maxPages = resp.Pagination.TotalPages
			log.Printf("maxPages for company search: %d", maxPages)
		}

		if len(resp.Accounts) == 0 {
			log.Println("No more companies to add")
			return nil
		}

		for _, fetchedCompany := range resp.Accounts {
			log.Println("Checking company:", fetchedCompany.Name, fetchedCompany.OrganizationID)
			if !fetchedCompany.CheckBasicFields() {
				log.Println("Skipping this company, failing basic fields check...")
				continue
			}

			exists, err := companies.ExistsByApolloID(ctx, fetchedCompany.OrganizationID, companies.StatusActive)
			if err != nil {
				log.Println("err:", err)
				return err
			}

			if exists {
				log.Printf("Skipping company: %s, %s because already exists", fetchedCompany.Name, fetchedCompany.OrganizationID)
				continue
			}

			dbCompanyID, err := companies.InsertIfNotExists(ctx, companies.CompanyForSet{
				ID:             uuid.New().String(),
				ApolloID:       fetchedCompany.OrganizationID,
				Status:         companies.StatusActive,
				Name:           fetchedCompany.Name,
				Website:        fetchedCompany.WebsiteURL,
				Industry:       "TECH",
				SubIndustry:    "",
				TechDetails:    "",
				CompanyDetails: "",
				Metadata:       "{}",
			})

			err = EnrichContactsForOrganizationID(ctx, fetchedCompany.OrganizationID, dbCompanyID, maxContactsPerCompany)
			if err != nil {
				if err == constants.ErrorNoPeopleFound {
					log.Println("no contacts found for the company, deleting it")
					err = companies.DeleteByID(ctx, dbCompanyID)
					if err != nil {
						log.Println("err:", err)
						return err
					}
					continue
				}
				log.Println("err:", err)
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
			log.Println("err:", err)
			return err
		}

		if pageNum == 1 {
			if len(resp.People) == 0 {
				err = constants.ErrorNoPeopleFound
				log.Println("warn err:", err)
				return err
			}
			maxPages = resp.Pagination.TotalPages
			log.Printf("maxPages for this company's person search: %d", maxPages)
		}

		if len(resp.People) == 0 {
			log.Println("No more persons to add")
			return nil
		}

		for _, fetchedPerson := range resp.People {
			log.Println("checking person:", fetchedPerson.Name)
			if !fetchedPerson.CheckBasicFields() {
				log.Println("Skipping this person, failing basic fields check...")
				continue
			}

			if fetchedPerson.Email == constants.ApolloLockedEmail {
				log.Println("enriching email")
				enrichContactResp, err := client.EnrichContactByApolloID(apollo.PersonMatchRequest{
					ID:                   fetchedPerson.ID,
					RevealPersonalEmails: true,
				})
				if err != nil {
					log.Println("err:", err)
					return err
				}
				log.Println("enriched email:", enrichContactResp.Person.Email)
				if enrichContactResp.Person.Email == "" {
					log.Println("enriched email is empty, skipping this person")
					continue
				}
				fetchedPerson.Email = enrichContactResp.Person.Email
			}

			exists, err := contacts.ExistsByApolloID(ctx, fetchedPerson.ID, contacts.StatusActive)
			if err != nil {
				log.Println("err:", err)
				return err
			}

			if exists {
				log.Printf("contact already exists, skipping...")
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
				log.Println("err:", err)
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

func GenerateAndSendEmails(ctx context.Context, count int, status string) (err error) {
	// Get Gmail token
	token, err := gmailtokens.GetLatestToken(ctx)
	if err != nil {
		log.Println("error:", err)
		return err
	}

	// Step 2 - Fetch contacts with company info
	contactsList, err := contacts.GetContactsWithCompanyInfo(ctx, count, status, "")
	if err != nil {
		log.Printf("ERROR: Failed to fetch contacts: %v", err)
		return err
	}

	log.Printf("Found %d contacts to send emails to", len(contactsList))

	// Get profile info for email generation
	profile, err := profileinfo.GetLatestActive(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to get profile info: %v", err)
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
		log.Printf("Processing contact: %s (%s) at %s", contact.ContactName, contact.ContactEmail, contact.CompanyName)

		// Step 3: Generate personalized email
		emailData := openai.EmailGenerationData{
			ContactName:       contact.ContactName,
			ContactRole:       contact.ContactRole,
			ContactLinkedIn:   contact.ContactLinkedIn,
			CompanyName:       contact.CompanyName,
			CompanyWebsite:    contact.CompanyWebsite,
			CompanyIndustry:   contact.CompanyIndustry,
			ProfileName:       profile.Name,
			ProfileExperience: profile.Experience,
			ProfileSkills:     profile.Skills,
			ProfileSummary:    profile.Summary,
		}

		subject, body, err := openaiClient.GeneratePersonalizedEmail(emailData)
		if err != nil {
			log.Printf("ERROR: Failed to generate email for %s: %v", contact.ContactEmail, err)
			emailsFailed++
			continue
		}

		emailsGenerated++
		log.Printf("Generated email for %s: %s", contact.ContactEmail, subject)

		// Log GENERATED stage
		metadata := map[string]interface{}{
			"contact_name": contact.ContactName,
			"contact_role": contact.ContactRole,
			"company_name": contact.CompanyName,
			"generated_at": time.Now().Format(time.RFC3339),
		}

		if err := email_logs.LogGenerated(ctx, contact.ContactID, contact.CompanyID, subject, body, metadata); err != nil {
			log.Printf("ERROR: Failed to log GENERATED stage for %s: %v", contact.ContactEmail, err)
		}

		// Step 4: Send email with resume
		if profile.ResumePath == "" {
			log.Printf("WARNING: No resume path found for profile, sending without attachment")
			err = gmail.SendSingleEmail(ctx, token.AccessToken, from, contact.ContactEmail, subject, body)
		} else {
			err = gmail.SendEmailWithAttachment(ctx, token.AccessToken, from, contact.ContactEmail, subject, body, profile.ResumePath)
		}

		if err != nil {
			log.Printf("ERROR: Failed to send email to %s: %v", contact.ContactEmail, err)
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
				log.Printf("ERROR: Failed to log ERROR stage for %s: %v", contact.ContactEmail, logErr)
			}
			continue
		}

		emailsSent++
		log.Printf("SUCCESS: Sent email to %s (%s)", contact.ContactEmail, subject)

		// Log SENT stage
		sentMetadata := map[string]interface{}{
			"contact_name":    contact.ContactName,
			"contact_role":    contact.ContactRole,
			"company_name":    contact.CompanyName,
			"sent_at":         time.Now().Format(time.RFC3339),
			"resume_attached": profile.ResumePath != "",
		}

		if err := email_logs.LogSent(ctx, contact.ContactID, contact.CompanyID, subject, body, sentMetadata); err != nil {
			log.Printf("ERROR: Failed to log SENT stage for %s: %v", contact.ContactEmail, err)
		}

		// Small delay to avoid rate limiting
		time.Sleep(2 * time.Second)
	}

	log.Printf("Email sending completed: %d generated, %d sent, %d failed", emailsGenerated, emailsSent, emailsFailed)

	return nil
}

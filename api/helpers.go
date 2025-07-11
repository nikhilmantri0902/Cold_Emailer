package api

import (
	"cold_emailer/apollo"
	"cold_emailer/constants"
	"cold_emailer/dbmodels/companies"
	"cold_emailer/dbmodels/contacts"
	"context"
	"log"
	"math"

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

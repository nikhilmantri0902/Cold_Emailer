package apollo

import (
	"cold_emailer/constants"
	"fmt"
)

// EnrichedCompany holds a company and its contacts
// This is for orchestration, not DB
type EnrichedCompany struct {
	Company  Company
	Contacts []Contact
}

// FetchEnrichedCompanies fetches up to 10 companies and their suitable contacts from Apollo
func FetchEnrichedCompanies(client *ApolloClient, maxCompanies int) ([]EnrichedCompany, error) {
	var enriched []EnrichedCompany
	companiesFetched := 0
	perCountry := maxCompanies / len(constants.TargetCountries)
	if perCountry == 0 {
		perCountry = 1
	}

	for _, country := range constants.TargetCountries {
		if companiesFetched >= maxCompanies {
			break
		}
		params := CompanySearchParams{
			Country:  country,
			Industry: "TECH",
			Page:     1,
			PerPage:  perCountry,
		}
		companies, err := client.SearchCompanies(params)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch companies for %s: %w", country, err)
		}
		for _, company := range companies {
			if companiesFetched >= maxCompanies {
				break
			}
			contacts, err := client.SearchContacts(company.ID)
			if err != nil {
				// Log and continue
				contacts = nil
			}
			enriched = append(enriched, EnrichedCompany{
				Company:  company,
				Contacts: contacts,
			})
			companiesFetched++
		}
	}
	return enriched, nil
}

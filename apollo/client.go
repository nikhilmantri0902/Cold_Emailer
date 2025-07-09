package apollo

import (
	"cold_emailer/constants"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ApolloClient handles requests to the Apollo API
// Reads API key from environment variable APOLLO_API_KEY
type ApolloClient struct {
	APIKey string
}

func NewApolloClient() *ApolloClient {
	apiKey := constants.APOLLO_API_KEY
	return &ApolloClient{APIKey: apiKey}
}

// CompanySearchParams defines parameters for company search
// (expand as needed)
type CompanySearchParams struct {
	Country  string
	Industry string
	Page     int
	PerPage  int
}

// Company represents a company from Apollo API
// (expand as needed)
type Company struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Website  string                 `json:"website"`
	Industry string                 `json:"industry"`
	Metadata map[string]interface{} `json:"-"` // Store raw Apollo data
}

// Contact represents a contact from Apollo API
// (expand as needed)
type Contact struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Email    string                 `json:"email"`
	Role     string                 `json:"title"`
	LinkedIn string                 `json:"linkedin_url"`
	Metadata map[string]interface{} `json:"-"`
}

// SearchCompanies fetches companies from Apollo API (stub)
func (c *ApolloClient) SearchCompanies(params CompanySearchParams) ([]Company, error) {
	baseURL := "https://api.apollo.io/v1/organizations/search"
	query := url.Values{}
	query.Set("api_key", c.APIKey)
	query.Set("q_organization_countries", params.Country)
	if params.Industry != "" {
		query.Set("q_organization_industries", params.Industry)
	}
	query.Set("page", fmt.Sprintf("%d", params.Page))
	query.Set("per_page", fmt.Sprintf("%d", params.PerPage))

	fullURL := fmt.Sprintf("%s?%s", baseURL, query.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call Apollo company search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Apollo company search failed: %s", resp.Status)
	}

	var result struct {
		Organizations []map[string]interface{} `json:"organizations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Apollo company search response: %w", err)
	}

	companies := make([]Company, 0, len(result.Organizations))
	for _, org := range result.Organizations {
		company := Company{
			Metadata: org,
		}
		if id, ok := org["id"].(string); ok {
			company.ID = id
		}
		if name, ok := org["name"].(string); ok {
			company.Name = name
		}
		if website, ok := org["website_url"].(string); ok {
			company.Website = website
		}
		if industry, ok := org["industry"].(string); ok {
			company.Industry = industry
		}
		companies = append(companies, company)
	}
	return companies, nil
}

// SuitableRoles is a helper for filtering contacts
var SuitableRoles = []string{"HR", "Recruiter", "CTO", "Tech Lead", "Technical Recruiter", "Talent Acquisition"}

func (c *ApolloClient) SearchContacts(companyID string) ([]Contact, error) {
	baseURL := "https://api.apollo.io/v1/people/search"
	query := url.Values{}
	query.Set("api_key", c.APIKey)
	query.Set("q_organization_ids", companyID)
	query.Set("per_page", "10")
	query.Set("page", "1")
	query.Set("q_titles", strings.Join(SuitableRoles, ","))

	fullURL := fmt.Sprintf("%s?%s", baseURL, query.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call Apollo contact search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Apollo contact search failed: %s", resp.Status)
	}

	var result struct {
		People []map[string]interface{} `json:"people"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Apollo contact search response: %w", err)
	}

	contacts := make([]Contact, 0, len(result.People))
	for _, p := range result.People {
		contact := Contact{
			Metadata: p,
		}
		if id, ok := p["id"].(string); ok {
			contact.ID = id
		}
		if name, ok := p["name"].(string); ok {
			contact.Name = name
		}
		if email, ok := p["email"].(string); ok {
			contact.Email = email
		}
		if title, ok := p["title"].(string); ok {
			contact.Role = title
		}
		if linkedin, ok := p["linkedin_url"].(string); ok {
			contact.LinkedIn = linkedin
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

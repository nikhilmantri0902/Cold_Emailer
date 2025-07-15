package apollo

import (
	"bytes"
	"cold_emailer/constants"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

// SearchContacts uses the POST /api/v1/mixed_people/search endpoint
func (c *ApolloClient) SearchContacts(req MixedPeopleSearchRequest) (mixedPeopleSearchResp MixedPeopleSearchResp, err error) {
	url := "https://api.apollo.io/api/v1/mixed_people/search"
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("accept", "application/json")
	httpReq.Header.Set("x-api-key", c.APIKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("error: %v", err)
		return mixedPeopleSearchResp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		log.Printf("Apollo error: %s", string(b))
		return mixedPeopleSearchResp, fmt.Errorf("Apollo contact search failed: %s", resp.Status)
	}

	var result MixedPeopleSearchResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("error: %v", err)
		return mixedPeopleSearchResp, err
	}

	return result, nil
}

// EnrichContactByApolloID uses the POST /api/v1/people/match endpoint
func (c *ApolloClient) EnrichContactByApolloID(req PersonMatchRequest) (personMatchResp PersonMatchResponse, err error) {
	url := "https://api.apollo.io/api/v1/people/match"
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("accept", "application/json")
	httpReq.Header.Set("x-api-key", c.APIKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("error: %v", err)
		return personMatchResp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		log.Printf("Apollo error: %s", string(b))
		return personMatchResp, fmt.Errorf("Apollo person match failed: %s", resp.Status)
	}

	var result PersonMatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("error: %v", err)
		return personMatchResp, err
	}

	return result, nil
}

// SearchCompaniesMixed uses the new POST /api/v1/mixed_companies/search endpoint
func (c *ApolloClient) SearchCompaniesMixed(req MixedCompanySearchRequest) (mixedCompanyResp MixedCompanySearchResp, err error) {
	url := "https://api.apollo.io/api/v1/mixed_companies/search"
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("accept", "application/json")
	httpReq.Header.Set("x-api-key", c.APIKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Println("error:", err)
		return mixedCompanyResp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		log.Printf("Apollo error: %s", string(b))
		return mixedCompanyResp, fmt.Errorf("Apollo company search failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&mixedCompanyResp); err != nil {
		log.Println("error:", err)
		return mixedCompanyResp, err
	}

	return mixedCompanyResp, nil
}

// GetOrganizationDetails fetches organization details by ID from Apollo
func (c *ApolloClient) GetOrganizationDetails(organizationID string) (ApolloOrgDetailResponse, error) {
	var respStruct ApolloOrgDetailResponse
	url := fmt.Sprintf("https://api.apollo.io/api/v1/organizations/%s", organizationID)
	httpReq, _ := http.NewRequest("GET", url, nil)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("accept", "application/json")
	httpReq.Header.Set("x-api-key", c.APIKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("error: %v", err)
		return respStruct, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		log.Printf("Apollo error: %s", string(b))
		return respStruct, fmt.Errorf("Apollo organization fetch failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		log.Printf("error: %v", err)
		return respStruct, err
	}

	return respStruct, nil
}

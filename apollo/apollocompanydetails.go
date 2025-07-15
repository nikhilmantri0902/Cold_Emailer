package apollo

import "strings"

// ApolloOrgDetailResponse is the response for GET /api/v1/organizations/{id}
type ApolloOrgDetailResponse struct {
	Organization ApolloOrgDetail `json:"organization"`
}

type ApolloOrgDetail struct {
	ID                           string                       `json:"id"`
	Name                         string                       `json:"name"`
	WebsiteURL                   string                       `json:"website_url"`
	BlogURL                      string                       `json:"blog_url"`
	AngellistURL                 string                       `json:"angellist_url"`
	LinkedinURL                  string                       `json:"linkedin_url"`
	TwitterURL                   string                       `json:"twitter_url"`
	FacebookURL                  string                       `json:"facebook_url"`
	PrimaryPhone                 PrimaryPhone                 `json:"primary_phone"`
	Languages                    []string                     `json:"languages"`
	AlexaRanking                 int                          `json:"alexa_ranking"`
	Phone                        string                       `json:"phone"`
	LinkedinUID                  string                       `json:"linkedin_uid"`
	FoundedYear                  int                          `json:"founded_year"`
	PubliclyTradedSymbol         string                       `json:"publicly_traded_symbol"`
	PubliclyTradedExchange       string                       `json:"publicly_traded_exchange"`
	LogoURL                      string                       `json:"logo_url"`
	CrunchbaseURL                string                       `json:"crunchbase_url"`
	PrimaryDomain                string                       `json:"primary_domain"`
	Industry                     string                       `json:"industry"`
	EstimatedNumEmployees        int                          `json:"estimated_num_employees"`
	Keywords                     []string                     `json:"keywords"`
	Industries                   []string                     `json:"industries"`
	SecondaryIndustries          []string                     `json:"secondary_industries"`
	SnippetsLoaded               bool                         `json:"snippets_loaded"`
	IndustryTagID                string                       `json:"industry_tag_id"`
	IndustryTagHash              map[string]string            `json:"industry_tag_hash"`
	RetailLocationCount          int                          `json:"retail_location_count"`
	RawAddress                   string                       `json:"raw_address"`
	StreetAddress                string                       `json:"street_address"`
	City                         string                       `json:"city"`
	State                        string                       `json:"state"`
	PostalCode                   string                       `json:"postal_code"`
	Country                      string                       `json:"country"`
	OwnedByOrganizationID        string                       `json:"owned_by_organization_id"`
	ShortDescription             string                       `json:"short_description"`
	Suborganizations             []interface{}                `json:"suborganizations"`
	NumSuborganizations          int                          `json:"num_suborganizations"`
	AnnualRevenuePrinted         string                       `json:"annual_revenue_printed"`
	AnnualRevenue                float64                      `json:"annual_revenue"`
	TotalFunding                 float64                      `json:"total_funding"`
	TotalFundingPrinted          string                       `json:"total_funding_printed"`
	LatestFundingRoundDate       string                       `json:"latest_funding_round_date"`
	LatestFundingStage           string                       `json:"latest_funding_stage"`
	FundingEvents                []ApolloOrgFundingEvent      `json:"funding_events"`
	TechnologyNames              []string                     `json:"technology_names"`
	CurrentTechnologies          []ApolloOrgCurrentTechnology `json:"current_technologies"`
	OrgChartRootPeopleIDs        []string                     `json:"org_chart_root_people_ids"`
	OrgChartSector               string                       `json:"org_chart_sector"`
	OrgChartRemoved              bool                         `json:"org_chart_removed"`
	OrgChartShowDepartmentFilter bool                         `json:"org_chart_show_department_filter"`
	AccountID                    string                       `json:"account_id"`
	Account                      ApolloOrgAccount             `json:"account"`
	HasIntentSignalAccount       bool                         `json:"has_intent_signal_account"`
	IntentSignalAccount          interface{}                  `json:"intent_signal_account"`
	GenericOrgInsights           interface{}                  `json:"generic_org_insights"`
	ShowIntent                   bool                         `json:"show_intent"`
	DetailViewLoaded             bool                         `json:"detail_view_loaded"`
}

func (o *ApolloOrgDetail) StringifyTechnologyArray() string {
	if len(o.TechnologyNames) == 0 {
		return ""
	}
	return strings.Join(o.TechnologyNames, ",")
}

type ApolloOrgFundingEvent struct {
	ID        string `json:"id"`
	Date      string `json:"date"`
	NewsURL   string `json:"news_url"`
	Type      string `json:"type"`
	Investors string `json:"investors"`
	Currency  string `json:"currency"`
}

type ApolloOrgCurrentTechnology struct {
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type ApolloOrgAccount struct {
	RawAddress                                 string                 `json:"raw_address"`
	StreetAddress                              string                 `json:"street_address"`
	City                                       string                 `json:"city"`
	State                                      string                 `json:"state"`
	Country                                    string                 `json:"country"`
	PostalCode                                 string                 `json:"postal_code"`
	ID                                         string                 `json:"id"`
	Domain                                     string                 `json:"domain"`
	Name                                       string                 `json:"name"`
	TeamID                                     string                 `json:"team_id"`
	OrganizationID                             string                 `json:"organization_id"`
	AccountStageID                             string                 `json:"account_stage_id"`
	Source                                     string                 `json:"source"`
	OriginalSource                             string                 `json:"original_source"`
	CreatorID                                  string                 `json:"creator_id"`
	OwnerID                                    string                 `json:"owner_id"`
	CreatedAt                                  string                 `json:"created_at"`
	Phone                                      string                 `json:"phone"`
	PhoneStatus                                string                 `json:"phone_status"`
	HubspotID                                  string                 `json:"hubspot_id"`
	SalesforceID                               string                 `json:"salesforce_id"`
	CrmOwnerID                                 string                 `json:"crm_owner_id"`
	ParentAccountID                            string                 `json:"parent_account_id"`
	LinkedinURL                                string                 `json:"linkedin_url"`
	AccountPlaybookStatuses                    []interface{}          `json:"account_playbook_statuses"`
	ExistenceLevel                             string                 `json:"existence_level"`
	LabelIDs                                   []string               `json:"label_ids"`
	TypedCustomFields                          map[string]interface{} `json:"typed_custom_fields"`
	CustomFieldErrors                          map[string]interface{} `json:"custom_field_errors"`
	Modality                                   string                 `json:"modality"`
	SourceDisplayName                          string                 `json:"source_display_name"`
	CrmRecordURL                               string                 `json:"crm_record_url"`
	TwitterURL                                 string                 `json:"twitter_url"`
	FacebookURL                                string                 `json:"facebook_url"`
	IntentStrength                             float64                `json:"intent_strength"`
	ShowIntent                                 bool                   `json:"show_intent"`
	OrganizationHeadcountSixMonthGrowth        float64                `json:"organization_headcount_six_month_growth"`
	OrganizationHeadcountTwelveMonthGrowth     float64                `json:"organization_headcount_twelve_month_growth"`
	OrganizationHeadcountTwentyFourMonthGrowth float64                `json:"organization_headcount_twenty_four_month_growth"`
	GenericOrgInsights                         interface{}            `json:"generic_org_insights"`
}

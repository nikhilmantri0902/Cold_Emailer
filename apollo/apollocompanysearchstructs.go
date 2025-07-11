package apollo

// MixedCompanySearchRequest is the request body for the new Apollo mixed companies search
type MixedCompanySearchRequest struct {
	OrganizationNumEmployeesRanges []string `json:"organization_num_employees_ranges,omitempty"`
	OrganizationLocations          []string `json:"organization_locations,omitempty"`
	QOrganizationKeywordTags       []string `json:"q_organization_keyword_tags,omitempty"`
	Page                           int      `json:"page"`
	PerPage                        int      `json:"per_page"`
}

// MixedCompanySearchRequest is the response body for the new Apollo mixed companies search
type MixedCompanySearchResp struct {
	Breadcrumbs          []Breadcrumb  `json:"breadcrumbs"`
	PartialResultsOnly   bool          `json:"partial_results_only"`
	HasJoin              bool          `json:"has_join"`
	DisableEUProspecting bool          `json:"disable_eu_prospecting"`
	PartialResultsLimit  int           `json:"partial_results_limit"`
	Pagination           Pagination    `json:"pagination"`
	Accounts             []Account     `json:"accounts"`
	Organizations        []interface{} `json:"organizations"`
	ModelIDs             []string      `json:"model_ids"`
	NumFetchResult       *int          `json:"num_fetch_result"` // can be null
	DerivedParams        DerivedParams `json:"derived_params"`
}

type Breadcrumb struct {
	Label           string `json:"label"`
	SignalFieldName string `json:"signal_field_name"`
	Value           string `json:"value"`
	DisplayName     string `json:"display_name"`
}

type Pagination struct {
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
	TotalEntries int `json:"total_entries"`
	TotalPages   int `json:"total_pages"`
}

type Account struct {
	ID                                         string                 `json:"id"`
	Name                                       string                 `json:"name"`
	WebsiteURL                                 string                 `json:"website_url"`
	BlogURL                                    *string                `json:"blog_url"`
	AngellistURL                               *string                `json:"angellist_url"`
	LinkedinURL                                string                 `json:"linkedin_url"`
	TwitterURL                                 string                 `json:"twitter_url"`
	FacebookURL                                string                 `json:"facebook_url"`
	PrimaryPhone                               *PrimaryPhone          `json:"primary_phone"`
	Languages                                  []string               `json:"languages"`
	AlexaRanking                               *int                   `json:"alexa_ranking"`
	Phone                                      *string                `json:"phone"`
	LinkedinUID                                string                 `json:"linkedin_uid"`
	FoundedYear                                int                    `json:"founded_year"`
	PubliclyTradedSymbol                       *string                `json:"publicly_traded_symbol"`
	PubliclyTradedExchange                     *string                `json:"publicly_traded_exchange"`
	LogoURL                                    string                 `json:"logo_url"`
	CrunchbaseURL                              *string                `json:"crunchbase_url"`
	PrimaryDomain                              string                 `json:"primary_domain"`
	SanitizedPhone                             *string                `json:"sanitized_phone"`
	OwnedByOrganizationID                      *string                `json:"owned_by_organization_id"`
	OrganizationRevenuePrinted                 *string                `json:"organization_revenue_printed"`
	OrganizationRevenue                        float64                `json:"organization_revenue"`
	OrganizationRawAddress                     string                 `json:"organization_raw_address"`
	OrganizationPostalCode                     *string                `json:"organization_postal_code"`
	OrganizationStreetAddress                  string                 `json:"organization_street_address"`
	OrganizationCity                           string                 `json:"organization_city"`
	OrganizationState                          string                 `json:"organization_state"`
	OrganizationCountry                        string                 `json:"organization_country"`
	SuggestLocationEnrichment                  bool                   `json:"suggest_location_enrichment"`
	RawAddress                                 string                 `json:"raw_address"`
	StreetAddress                              string                 `json:"street_address"`
	City                                       string                 `json:"city"`
	State                                      string                 `json:"state"`
	Country                                    string                 `json:"country"`
	PostalCode                                 *string                `json:"postal_code"`
	Domain                                     string                 `json:"domain"`
	TeamID                                     string                 `json:"team_id"`
	OrganizationID                             string                 `json:"organization_id"`
	AccountStageID                             string                 `json:"account_stage_id"`
	Source                                     string                 `json:"source"`
	OriginalSource                             string                 `json:"original_source"`
	CreatorID                                  string                 `json:"creator_id"`
	OwnerID                                    *string                `json:"owner_id"`
	CreatedAt                                  string                 `json:"created_at"`
	PhoneStatus                                string                 `json:"phone_status"`
	HubspotID                                  *string                `json:"hubspot_id"`
	SalesforceID                               *string                `json:"salesforce_id"`
	CrmOwnerID                                 *string                `json:"crm_owner_id"`
	ParentAccountID                            *string                `json:"parent_account_id"`
	AccountPlaybookStatuses                    []interface{}          `json:"account_playbook_statuses"`
	ExistenceLevel                             string                 `json:"existence_level"`
	LabelIDs                                   []string               `json:"label_ids"`
	TypedCustomFields                          map[string]interface{} `json:"typed_custom_fields"`
	CustomFieldErrors                          map[string]interface{} `json:"custom_field_errors"`
	Modality                                   string                 `json:"modality"`
	SourceDisplayName                          string                 `json:"source_display_name"`
	CrmRecordURL                               *string                `json:"crm_record_url"`
	ContactEmailerCampaignIDs                  []interface{}          `json:"contact_emailer_campaign_ids"`
	ContactCampaignStatusTally                 map[string]interface{} `json:"contact_campaign_status_tally"`
	NumContacts                                int                    `json:"num_contacts"`
	LastActivityDate                           *string                `json:"last_activity_date"`
	IntentStrength                             *float64               `json:"intent_strength"`
	ShowIntent                                 bool                   `json:"show_intent"`
	OrganizationHeadcountSixMonthGrowth        float64                `json:"organization_headcount_six_month_growth"`
	OrganizationHeadcountTwelveMonthGrowth     float64                `json:"organization_headcount_twelve_month_growth"`
	OrganizationHeadcountTwentyFourMonthGrowth float64                `json:"organization_headcount_twenty_four_month_growth"`
}

// CheckBasicFields checks if at least a basic set of fields exist
func (a *Account) CheckBasicFields() bool {
	if a.OrganizationID != "" && a.Name != "" && a.WebsiteURL != "" {
		return true
	}
	return false
}

type PrimaryPhone struct {
	Number          string `json:"number"`
	Source          string `json:"source"`
	SanitizedNumber string `json:"sanitized_number"`
}

type DerivedParams struct {
	RecommendationConfigID string `json:"recommendation_config_id"`
}

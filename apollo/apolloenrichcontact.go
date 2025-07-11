package apollo

// PersonMatchRequest is the request body for Apollo person match endpoint
type PersonMatchRequest struct {
	ID                   string `json:"id"`
	RevealPersonalEmails bool   `json:"reveal_personal_emails"`
}

// EmploymentHistory represents employment history in the person match response
type PersonEmploymentHistory struct {
	ID               string  `json:"id"`
	CreatedAt        *string `json:"created_at"`
	Current          bool    `json:"current"`
	Degree           *string `json:"degree"`
	Description      *string `json:"description"`
	Emails           *string `json:"emails"`
	EndDate          *string `json:"end_date"`
	GradeLevel       *string `json:"grade_level"`
	Kind             *string `json:"kind"`
	Major            *string `json:"major"`
	OrganizationID   *string `json:"organization_id"`
	OrganizationName *string `json:"organization_name"`
	RawAddress       *string `json:"raw_address"`
	StartDate        *string `json:"start_date"`
	Title            *string `json:"title"`
	UpdatedAt        *string `json:"updated_at"`
	Key              string  `json:"key"`
}

// Organization represents the organization in the person match response
type PersonOrganization struct {
	ID                     string            `json:"id"`
	Name                   string            `json:"name"`
	WebsiteURL             string            `json:"website_url"`
	BlogURL                *string           `json:"blog_url"`
	AngellistURL           *string           `json:"angellist_url"`
	LinkedinURL            string            `json:"linkedin_url"`
	TwitterURL             *string           `json:"twitter_url"`
	FacebookURL            *string           `json:"facebook_url"`
	PrimaryPhone           *PrimaryPhone     `json:"primary_phone"`
	Languages              []string          `json:"languages"`
	AlexaRanking           int               `json:"alexa_ranking"`
	Phone                  string            `json:"phone"`
	LinkedinUID            string            `json:"linkedin_uid"`
	FoundedYear            int               `json:"founded_year"`
	PubliclyTradedSymbol   *string           `json:"publicly_traded_symbol"`
	PubliclyTradedExchange *string           `json:"publicly_traded_exchange"`
	LogoURL                string            `json:"logo_url"`
	CrunchbaseURL          *string           `json:"crunchbase_url"`
	PrimaryDomain          string            `json:"primary_domain"`
	SanitizedPhone         string            `json:"sanitized_phone"`
	Industry               string            `json:"industry"`
	EstimatedNumEmployees  int               `json:"estimated_num_employees"`
	Keywords               []string          `json:"keywords"`
	Industries             []string          `json:"industries"`
	SecondaryIndustries    []string          `json:"secondary_industries"`
	SnippetsLoaded         bool              `json:"snippets_loaded"`
	IndustryTagID          string            `json:"industry_tag_id"`
	IndustryTagHash        map[string]string `json:"industry_tag_hash"`
	RetailLocationCount    int               `json:"retail_location_count"`
	RawAddress             string            `json:"raw_address"`
	StreetAddress          string            `json:"street_address"`
	City                   string            `json:"city"`
	State                  string            `json:"state"`
	PostalCode             string            `json:"postal_code"`
	Country                string            `json:"country"`
	OwnedByOrganizationID  *string           `json:"owned_by_organization_id"`
	ShortDescription       string            `json:"short_description"`
	Suborganizations       []interface{}     `json:"suborganizations"`
	NumSuborganizations    int               `json:"num_suborganizations"`
	AnnualRevenuePrinted   string            `json:"annual_revenue_printed"`
	AnnualRevenue          float64           `json:"annual_revenue"`
	FundingEvents          []interface{}     `json:"funding_events"`
	TechnologyNames        []string          `json:"technology_names"`
	CurrentTechnologies    []interface{}     `json:"current_technologies"`
	OrgChartRootPeopleIDs  []string          `json:"org_chart_root_people_ids"`
	OrgChartSector         string            `json:"org_chart_sector"`
}

// Person represents a person in the person match response
type Person struct {
	ID                          string                    `json:"id"`
	FirstName                   string                    `json:"first_name"`
	LastName                    string                    `json:"last_name"`
	Name                        string                    `json:"name"`
	LinkedinURL                 string                    `json:"linkedin_url"`
	Title                       string                    `json:"title"`
	EmailStatus                 string                    `json:"email_status"`
	PhotoURL                    string                    `json:"photo_url"`
	TwitterURL                  *string                   `json:"twitter_url"`
	GithubURL                   *string                   `json:"github_url"`
	FacebookURL                 *string                   `json:"facebook_url"`
	ExtrapolatedEmailConfidence float64                   `json:"extrapolated_email_confidence"`
	Headline                    string                    `json:"headline"`
	Email                       string                    `json:"email"`
	OrganizationID              string                    `json:"organization_id"`
	EmploymentHistory           []PersonEmploymentHistory `json:"employment_history"`
	State                       *string                   `json:"state"`
	City                        *string                   `json:"city"`
	Country                     string                    `json:"country"`
	Organization                PersonOrganization        `json:"organization"`
	ShowIntent                  bool                      `json:"show_intent"`
	EmailDomainCatchall         bool                      `json:"email_domain_catchall"`
	RevealedForCurrentTeam      bool                      `json:"revealed_for_current_team"`
	PersonalEmails              []interface{}             `json:"personal_emails"`
	Departments                 []string                  `json:"departments"`
	Subdepartments              []string                  `json:"subdepartments"`
	Functions                   []string                  `json:"functions"`
	Seniority                   string                    `json:"seniority"`
}

func (p *Person) CheckBasicFields() bool {
	if p.Name != "" && p.Title != "" && p.LinkedinURL != "" {
		return true
	}
	return false
}

// PersonMatchResponse is the response for Apollo person match endpoint
type PersonMatchResponse struct {
	Person Person `json:"person"`
}

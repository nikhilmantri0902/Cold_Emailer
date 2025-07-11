package apollo

// MixedPeopleSearchRequest is the request body for Apollo mixed people search
type MixedPeopleSearchRequest struct {
	PersonTitles    []string `json:"person_titles"`
	OrganizationIDs []string `json:"organization_ids"`
	Page            int      `json:"page"`
	PerPage         int      `json:"per_page"`
}

// ContactEmail represents an email object in the contact_emails array
type ContactEmail struct {
	Email                        string   `json:"email"`
	EmailMD5                     string   `json:"email_md5"`
	EmailSHA256                  string   `json:"email_sha256"`
	EmailStatus                  string   `json:"email_status"`
	ExtrapolatedEmailConfidence  float64  `json:"extrapolated_email_confidence"`
	Position                     int      `json:"position"`
	FreeDomain                   bool     `json:"free_domain"`
	Source                       string   `json:"source"`
	ThirdPartyVendorName         string   `json:"third_party_vendor_name"`
	VendorValidationStatuses     []string `json:"vendor_validation_statuses"`
	EmailNeedsTickling           bool     `json:"email_needs_tickling"`
	EmailTrueStatus              string   `json:"email_true_status"`
	EmailStatusUnavailableReason string   `json:"email_status_unavailable_reason"`
}

// PhoneNumber represents a phone number object in the phone_numbers array
type PhoneNumber struct {
	RawNumber                string                 `json:"raw_number"`
	SanitizedNumber          string                 `json:"sanitized_number"`
	Type                     string                 `json:"type"`
	Position                 int                    `json:"position"`
	Status                   string                 `json:"status"`
	DNCStatus                *string                `json:"dnc_status"`
	DNC_OtherInfo            map[string]interface{} `json:"dnc_other_info"`
	DialerFlags              *string                `json:"dialer_flags"`
	SourceName               *string                `json:"source_name"`
	VendorValidationStatuses []string               `json:"vendor_validation_statuses"`
	ThirdPartyVendorName     *string                `json:"third_party_vendor_name"`
}

// EmploymentHistory represents an employment history object
type EmploymentHistory struct {
	ID               string  `json:"id"`
	Current          bool    `json:"current"`
	OrganizationID   *string `json:"organization_id"`
	OrganizationName *string `json:"organization_name"`
	Title            *string `json:"title"`
	StartDate        *string `json:"start_date"`
	EndDate          *string `json:"end_date"`
}

// Account represents the account object in the contact
type ContactAccount struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WebsiteURL  string `json:"website_url"`
	LinkedInURL string `json:"linkedin_url"`
	TwitterURL  string `json:"twitter_url"`
	FacebookURL string `json:"facebook_url"`
	Phone       string `json:"phone"`
	Domain      string `json:"domain"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
}

// Organization represents the organization object in the contact
type Organization struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	WebsiteURL     string `json:"website_url"`
	LinkedInURL    string `json:"linkedin_url"`
	TwitterURL     string `json:"twitter_url"`
	FacebookURL    string `json:"facebook_url"`
	Phone          string `json:"phone"`
	SanitizedPhone string `json:"sanitized_phone"`
	Domain         string `json:"domain"`
}

// MixedPersonContact represents a contact in the mixed people search response
type MixedPersonContact struct {
	ID                string              `json:"id"`
	FirstName         string              `json:"first_name"`
	LastName          string              `json:"last_name"`
	Name              string              `json:"name"`
	LinkedInURL       string              `json:"linkedin_url"`
	Title             string              `json:"title"`
	PersonID          string              `json:"person_id"`
	Email             string              `json:"email"`
	OrganizationName  string              `json:"organization_name"`
	OrganizationID    string              `json:"organization_id"`
	Headline          string              `json:"headline"`
	City              string              `json:"city"`
	Country           string              `json:"country"`
	State             string              `json:"state"`
	EmailStatus       string              `json:"email_status"`
	AccountID         string              `json:"account_id"`
	SanitizedPhone    string              `json:"sanitized_phone"`
	CreatedAt         string              `json:"created_at"`
	UpdatedAt         string              `json:"updated_at"`
	PhotoURL          string              `json:"photo_url"`
	ContactEmails     []ContactEmail      `json:"contact_emails"`
	PhoneNumbers      []PhoneNumber       `json:"phone_numbers"`
	EmploymentHistory []EmploymentHistory `json:"employment_history"`
	Account           *ContactAccount     `json:"account"`
	Organization      *Organization       `json:"organization"`
}

type MixedPersonSearchPeople struct {
	ID                          string              `json:"id"`
	FirstName                   string              `json:"first_name"`
	LastName                    string              `json:"last_name"`
	Name                        string              `json:"name"`
	LinkedinURL                 string              `json:"linkedin_url"`
	Title                       string              `json:"title"`
	EmailStatus                 string              `json:"email_status"`
	PhotoURL                    string              `json:"photo_url"`
	TwitterURL                  *string             `json:"twitter_url"`
	GithubURL                   *string             `json:"github_url"`
	FacebookURL                 *string             `json:"facebook_url"`
	ExtrapolatedEmailConfidence float64             `json:"extrapolated_email_confidence"`
	Headline                    string              `json:"headline"`
	Email                       string              `json:"email"`
	OrganizationID              string              `json:"organization_id"`
	EmploymentHistory           []EmploymentHistory `json:"employment_history"`
	State                       string              `json:"state"`
	City                        string              `json:"city"`
	Country                     string              `json:"country"`
	Organization                *Organization       `json:"organization"`
	AccountID                   string              `json:"account_id"`
	Account                     *ContactAccount     `json:"account"`
	Departments                 []string            `json:"departments"`
	Subdepartments              []string            `json:"subdepartments"`
	Seniority                   string              `json:"seniority"`
	Functions                   []string            `json:"functions"`
	IntentStrength              *string             `json:"intent_strength"`
	ShowIntent                  bool                `json:"show_intent"`
	EmailDomainCatchall         bool                `json:"email_domain_catchall"`
	RevealedForCurrentTeam      bool                `json:"revealed_for_current_team"`
}

func (m *MixedPersonSearchPeople) CheckBasicFields() bool {
	if m.Name != "" && m.LinkedinURL != "" && m.Title != "" {
		return true
	}
	return false
}

// MixedPeopleSearchResp is the response for Apollo mixed people search
type MixedPeopleSearchResp struct {
	Breadcrumbs          []Breadcrumb              `json:"breadcrumbs"`
	PartialResultsOnly   bool                      `json:"partial_results_only"`
	HasJoin              bool                      `json:"has_join"`
	DisableEUProspecting bool                      `json:"disable_eu_prospecting"`
	PartialResultsLimit  int                       `json:"partial_results_limit"`
	Pagination           Pagination                `json:"pagination"`
	Contacts             []MixedPersonContact      `json:"contacts"`
	ModelIDs             []string                  `json:"model_ids"`
	People               []MixedPersonSearchPeople `json:"people"`
	NumFetchResult       *int                      `json:"num_fetch_result"` // can be null
	DerivedParams        DerivedParams             `json:"derived_params"`
}

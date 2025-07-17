package contacts

type Contact struct {
	ID          string `db:"id"`
	ApolloID    string `db:"apollo_id"`
	CreatedAt   string `db:"created_at"`
	CompanyID   string `db:"company_id"`
	Status      string `db:"status"`
	Name        string `db:"name"`
	EmailID     string `db:"email_id"`
	PhoneNumber string `db:"phone_number"`
	LinkedInURL string `db:"linkedin_url"`
	Role        string `db:"role"`
	Metadata    string `db:"metadata"`
}

type ContactForSet struct {
	ID          string
	ApolloID    string
	CompanyID   string
	Status      string
	Name        string
	EmailID     string
	PhoneNumber string
	LinkedInURL string
	Role        string
	Metadata    string
}

// ContactWithCompany represents a contact with company information
type ContactWithCompany struct {
	ContactID       string `db:"contact_id"`
	ContactName     string `db:"contact_name"`
	ContactEmail    string `db:"contact_email"`
	ContactRole     string `db:"contact_role"`
	ContactLinkedIn string `db:"contact_linkedin"`
	ContactPhone    string `db:"contact_phone"`
	ContactStatus   string `db:"contact_status"`
	CompanyID       string `db:"company_id"`
	CompanyName     string `db:"company_name"`
	CompanyWebsite  string `db:"company_website"`
	CompanyIndustry string `db:"company_industry"`
	CompanyTech     string `db:"tech_details"`
	CompanyDetails  string `db:"company_details"`
}

type FollowUpCandidate struct {
	ContactID       string
	CompanyID       string
	ContactName     string
	ContactEmail    string
	ContactRole     string
	ContactLinkedIn string
	ContactPhone    string
	ContactStatus   string
	CompanyName     string
	CompanyWebsite  string
	CompanyIndustry string
	CompanyTech     string
	CompanyDetails  string
	EmailSubject    string
	EmailBody       string
	Metadata        string
}

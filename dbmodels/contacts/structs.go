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

package companies

type Company struct {
	ID             string `db:"id"`
	ApolloID       string `db:"apollo_id"`
	CreatedAt      string `db:"created_at"`
	Status         string `db:"status"`
	Name           string `db:"name"`
	Website        string `db:"website"`
	Industry       string `db:"industry"`
	SubIndustry    string `db:"sub_industry"`
	TechDetails    string `db:"tech_details"`
	CompanyDetails string `db:"company_details"`
	Metadata       string `db:"metadata"`
}

type CompanyForSet struct {
	ID             string
	ApolloID       string
	Status         string
	Name           string
	Website        string
	Industry       string
	SubIndustry    string
	TechDetails    string
	CompanyDetails string
	Metadata       string
}

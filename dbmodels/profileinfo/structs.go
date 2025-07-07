package profileinfo

type StructForSet struct {
	ID          string `db:"id"`
	Status      string `db:"status"`
	Name        string `db:"name"`
	Email       string `db:"email"`
	Phone       string `db:"phone"`
	LinkedInURL string `db:"linkedin_url"`
	Experience  string `db:"experience"`
	Skills      string `db:"skills"`
	Summary     string `db:"summary"`
	ResumePath  string `db:"resume_path"`
	Metadata    string `db:"metadata"`
}

type Struct struct {
	ID          string `db:"id"`
	CreatedAt   string `db:"created_at"`
	Status      string `db:"status"`
	Name        string `db:"name"`
	Email       string `db:"email"`
	Phone       string `db:"phone"`
	LinkedInURL string `db:"linkedin_url"`
	Experience  string `db:"experience"`
	Skills      string `db:"skills"`
	Summary     string `db:"summary"`
	ResumePath  string `db:"resume_path"`
	Metadata    string `db:"metadata"`
}

# Cold Emailer

A Golang-based API-first project to automate personalized cold outreach to recruiters and CTOs, leveraging OpenAI for email generation and Gmail API for sending emails (with resume attachment).
Built for API usage (testable via Postman).

## Features (WIP)

* Upload your profile & resume
* Upload list of targets (CTOs, recruiters, companies)
* Generate personalized emails using OpenAI
* Send emails (Gmail API integration)
* Track status and logs

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/nikhilmantri0902/Cold_Emailer
cd Cold_Emailer
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Setup environment variables

Copy `.env.example` into `.env` and fill in your secrets:

```env
PORT=8080
OPENAI_API_KEY=your-openai-api-key
GMAIL_CLIENT_ID=your-gmail-client-id
GMAIL_CLIENT_SECRET=your-gmail-client-secret
GMAIL_REDIRECT_URI=http://localhost:8080/oauth2callback
DB_URL=coldemailer.db
```

### 4. Run the server

```bash
go run cmd/server/main.go
```

### 5. Test the health endpoint

Use curl or Postman:

```bash
curl http://localhost:8080/health
# Response: {"status":"ok"}
```

## API Endpoints

* `POST /api/profile` — Upload profile/resume
* `POST /api/targets` — Upload targets (recruiters, CTOs)
* `POST /api/generate-email` — Generate email draft (OpenAI)
* `POST /api/send-email` — Send email via Gmail
* `GET /api/status` — Email send status
* `GET /api/logs` — Logs

*(All endpoints are currently stubs / in development.)*

---

## License

MIT

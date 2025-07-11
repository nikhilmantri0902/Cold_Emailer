# Cold Emailer

A comprehensive Golang-based cold email outreach system that automates personalized outreach to recruiters and CTOs. The system integrates Apollo API for company and contact discovery, OpenAI for personalized email generation, and Gmail API for sending emails with resume attachments.

## Features

### ✅ Implemented Features
* **Profile Management**: Upload and store your profile information and resume
* **Apollo Integration**: Automatically discover companies and contacts using Apollo's API
* **Smart Enrichment**: Enrich database with companies and contacts filtered by country and employee range
* **AI-Powered Email Generation**: Generate personalized emails using OpenAI with contact and company context
* **Gmail Integration**: Send emails with resume attachments via Gmail API
* **Comprehensive Logging**: Track email generation, sending, and error states
* **Rate Limiting**: Built-in rate limiting for API calls
* **Database Management**: PostgreSQL with proper schema and migrations

### 🔄 Pipeline Flow
1. **Enrichment**: Trigger company/contact discovery from Apollo
2. **Email Generation**: Create personalized emails for selected contacts
3. **Email Sending**: Send emails with resume attachments via Gmail
4. **Logging**: Track every step in the email_logs table

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

Create a `.env` file with the following variables:

```env
PORT=8080
OPENAI_API_KEY=your-openai-api-key
GMAIL_CLIENT_ID=your-gmail-client-id
GMAIL_CLIENT_SECRET=your-gmail-client-secret
GMAIL_REDIRECT_URI=http://localhost:8080/oauth2callback
APOLLO_API_KEY=your-apollo-api-key
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=cold_emailer
```

### 4. Setup Database

The system uses PostgreSQL. You can run it using Docker:

```bash
docker-compose up -d postgres
```

The migrations will automatically run to create the required tables:
- `profile_info` - User profile and resume information
- `gmail_tokens` - Gmail OAuth tokens
- `companies` - Companies discovered from Apollo
- `contacts` - Contacts discovered from Apollo
- `email_logs` - Email generation and sending logs

### 5. Run the server

```bash
docker-compose up --build
```

Or run directly:

```bash
go run cmd/server/main.go
```

### 6. Test the health endpoint

```bash
curl http://localhost:8080/health
# Response: {"status":"ok"}
```

## API Endpoints

### Profile Management
* `POST /api/profile` — Upload profile information and resume
* `GET /api/profile` — Get current profile information

### Apollo Integration
* `POST /api/enrich` — Trigger company and contact discovery from Apollo
  - Automatically fetches companies and contacts by country
  - Stops when 10 new companies are added per country
  - Stores Apollo metadata for future reference

### Email Pipeline
* `POST /api/send-few-initial-emails` — Send personalized emails to 10 contacts
  - Generates personalized emails using OpenAI
  - Sends emails with resume attachments via Gmail
  - Logs each step (GENERATED, SENT, ERROR)
  - Excludes contacts that already have "SENT" emails

### Gmail OAuth
* `GET /api/auth/gmail` — Initiate Gmail OAuth flow
* `GET /api/oauth2callback` — Gmail OAuth callback handler

### Status & Logs
* `GET /api/health` — Health check endpoint
* Database queries available for email logs and status tracking

## Database Schema

### Companies Table
- Company information from Apollo
- Includes name, website, employee count, industry, etc.
- Apollo metadata for tracking

### Contacts Table
- Contact information from Apollo
- Includes name, email, title, company association
- Apollo metadata for tracking

### Email Logs Table
- Tracks email generation and sending status
- States: GENERATED, SENT, ERROR
- Includes error messages and timestamps

### Profile Info Table
- User profile information
- Resume file path for attachments

## Usage Workflow

1. **Setup Profile**: Upload your profile and resume
2. **Authenticate Gmail**: Complete Gmail OAuth flow
3. **Enrich Database**: Trigger Apollo enrichment to discover companies/contacts
4. **Send Emails**: Use the send-few-initial-emails endpoint to start outreach
5. **Monitor Logs**: Check email_logs table for status and results

## Technical Architecture

- **Backend**: Golang with Gin framework
- **Database**: PostgreSQL with proper migrations
- **External APIs**: 
  - Apollo API for company/contact discovery
  - OpenAI API for email generation
  - Gmail API for email sending
- **File Storage**: Local file system for resume storage
- **Authentication**: Gmail OAuth 2.0

## Rate Limiting

The system includes built-in rate limiting for:
- Apollo API calls (company and contact search)
- OpenAI API calls (email generation)
- Gmail API calls (email sending)

## Error Handling

Comprehensive error handling with:
- Database transaction rollbacks
- API call retries
- Detailed error logging
- Graceful degradation

---

## License

MIT

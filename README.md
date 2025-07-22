# Cold Emailer - AI-Powered Job Application Outreach System

A comprehensive Golang-based cold email outreach system that automates personalized job application emails to recruiters and CTOs. The system integrates Apollo API for company and contact discovery, OpenAI for intelligent email generation, and Gmail API for sending emails with resume attachments. Includes a web-based dashboard for monitoring and management.

## 🚀 Features

### ✅ Core Functionality
- **Profile Management**: Upload and store your profile information and resume
- **Apollo Integration**: Automatically discover companies and contacts using Apollo's API
- **Smart Company & Contact Enrichment**: Enrich database with companies and contacts filtered by country, employee range, and industry
- **AI-Powered Email Generation**: Generate personalized emails using OpenAI with contact and company context
- **Gmail Integration**: Send emails with resume attachments via Gmail API OAuth
- **Comprehensive Logging**: Track email generation, sending, and error states with detailed logs
- **Rate Limiting**: Built-in rate limiting for all external API calls
- **Database Management**: PostgreSQL with proper schema, migrations, and indexing

### 🖥 Frontend Dashboard
- **Web-Based Interface**: Clean, responsive dashboard for monitoring and management
- **Email Logs Viewer**: Paginated display of all email activities with filtering options
- **Company Browser**: Interactive company directory with search and filtering
- **Contact Management**: View detailed contact information organized by company
- **System Configuration**: Display of current targeting parameters and filters
- **Real-time Data**: Live updates from the backend API
- **Tab-Based Navigation**: Organized interface with separate sections for different functionalities

### 🔄 Complete Pipeline Flow
1. **Profile Setup**: Upload your profile information and resume
2. **Gmail Authentication**: Complete OAuth flow for email sending
3. **Database Enrichment**: Discover and store companies/contacts from Apollo API
4. **Company Details Backfill**: Enhance company information with detailed data
5. **Email Generation**: Create personalized emails using AI for selected contacts
6. **Email Sending**: Send emails with resume attachments and track status
7. **Follow-up Management**: Send follow-up emails to previously contacted leads

## 🛠 Technology Stack

- **Backend**: Go 1.23 with Gin web framework
- **Frontend**: Static HTML/CSS/JavaScript with Nginx
- **Database**: PostgreSQL 15 with migrations
- **External APIs**: 
  - Apollo API for company/contact discovery
  - OpenAI API (GPT-3.5-turbo) for email generation
  - Gmail API for email sending
- **File Storage**: Local file system with organized storage structure
- **Authentication**: Gmail OAuth 2.0
- **Web Server**: Nginx (reverse proxy + static file serving)
- **Deployment**: Docker & Docker Compose
- **Logging**: Structured logging with zerolog

## 📋 Prerequisites

- Go 1.23+
- PostgreSQL 15+
- Docker & Docker Compose (recommended)
- API Keys for:
  - OpenAI API
  - Apollo API
  - Gmail API (OAuth credentials)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/nikhilmantri0902/Cold_Emailer
cd Cold_Emailer
```

### 2. Environment Configuration

Create a `.env` file in the root directory:

```env
# Server Configuration
PORT=8000

# OpenAI Configuration
OPENAI_API_KEY=your-openai-api-key
OPENAI_MODEL=gpt-3.5-turbo
OPENAI_TEMPERATURE=0.7
OPENAI_MAX_COMPLETION_TOKENS=512

# Gmail OAuth Configuration
GMAIL_CLIENT_ID=your-gmail-client-id
GMAIL_CLIENT_SECRET=your-gmail-client-secret
GMAIL_REDIRECT_URI=http://localhost:8000/gmail-oauth2callback

# Apollo API Configuration
APOLLO_API_KEY=your-apollo-api-key

# Database Configuration (for local development)
DB_HOST=localhost
DB_PORT=5432
DB_USER=coldemailer
DB_PASSWORD=coldemailer
DB_NAME=coldemailer
```

### 3. Using Docker Compose (Recommended)

```bash
# Start all services (backend, frontend, database)
docker-compose up --build

# Run in background
docker-compose up -d --build

# Services started:
# - Backend API: http://localhost:8000
# - Frontend Dashboard: http://localhost:3000  
# - PostgreSQL Database: localhost:5432
```

### 4. Manual Setup (Alternative)

```bash
# Install dependencies
go mod tidy

# Start PostgreSQL (ensure it's running)
# Update constants/constants.go with your DB connection string

# Run migrations and start server
go run cmd/server/main.go
```

### 5. Verify Installation

```bash
# Check backend API
curl http://localhost:8000/health
# Response: {"status":"ok"}

# Access frontend dashboard
open http://localhost:3000
```

## 🌐 Access Points

### Frontend Dashboard
- **URL**: http://localhost:3000
- **Features**: 
  - Email Logs tab with pagination and filtering
  - Companies & Contacts tab with search functionality
  - System configuration display
  - Real-time data from backend API

### Backend API
- **URL**: http://localhost:8000
- **Health Check**: http://localhost:8000/health
- **API Documentation**: See API Endpoints section below

## 📡 API Endpoints

### Profile Management

#### Upload Profile & Resume
```http
POST /api/profile
Content-Type: multipart/form-data

Form Data:
- name: string (required)
- email: string (required) 
- phone: string (optional)
- linkedin_url: string (optional)
- experience: string (optional)
- skills: string (optional)
- summary: string (optional)
- resume: file (optional - PDF/DOC)
```



### Database Enrichment

#### Enrich Database with Companies & Contacts
```http
POST /api/enrich-database
Content-Type: application/json

{
  "count_new_companies": 50,
  "max_contacts_per_company": 5
}
```
- Fetches companies from Apollo API filtered by:
  - Target countries: Germany, Netherlands, Canada, Sweden, Finland, Norway, Ireland, UK, Luxembourg, UAE, Singapore, Australia
  - Employee ranges: 1-10, 11-50, 51-200, 201-500
  - Industries: Technology, Software, IT, AI
- Discovers contacts with suitable roles: CTO, Talent Acquisition, Technical Recruiter, Product Manager

#### Backfill Company Details
```http
POST /api/backfill-company-details
```
- Enhances existing companies with detailed information from Apollo

### Email Operations

#### Generate Single Email
```http
POST /api/generate-email
Content-Type: application/json

{
  "contact_id": "uuid-of-contact"
}
```

#### Send Single Email
```http
POST /api/send-single-email
Content-Type: application/json

{
  "contact_id": "uuid-of-contact"
}
```

#### Send Initial Outreach Emails (Batch)
```http
POST /api/send-few-initial-emails
```
- Automatically selects up to 10 contacts without prior email history
- Generates personalized emails using AI
- Sends emails with resume attachments
- Logs all activities with status tracking

#### Send Follow-up Emails (Batch)
```http
POST /api/send-few-follow-up-emails
```
- Sends follow-up emails to previously contacted leads
- Excludes recent contacts (within last 7 days)

### Gmail Authentication

#### Initiate Gmail OAuth
```http
GET /gmail-auth-initiate
```
- Redirects to Gmail OAuth consent screen

#### OAuth Callback (Automatic)
```http
GET /gmail-oauth2callback?code=...&state=...
```
- Handles OAuth callback and stores tokens

### Frontend API Endpoints

#### Get Email Logs (with Pagination)
```http
GET /api/email-logs?page=1&limit=10&stage=INITIAL
```
- Returns paginated email logs with company and contact details
- **Query Parameters**:
  - `page`: Page number (default: 1)
  - `limit`: Records per page (default: 10, max: 100)
  - `stage`: Filter by email stage (INITIAL, FOLLOWUP) - optional

#### Get All Companies
```http
GET /api/companies
```
- Returns all companies in the database with full details

#### Get Company Contacts
```http
GET /api/companies/{company_id}/contacts
```
- Returns all contacts for a specific company

#### Get System Configuration
```http
GET /api/config
```
- Returns current targeting configuration (countries, roles, company sizes, industries)

### Monitoring & Status

#### Health Check
```http
GET /health
```

#### System Status
```http
GET /api/status
```
- Returns database connection status and system health

#### Email Logs (Legacy)
```http
GET /api/logs
```
- Returns recent email generation and sending logs

## 🗄 Database Schema

### Tables Overview

#### `profile_info`
Stores user profile and resume information
```sql
- id (UUID, Primary Key)
- created_at (Timestamp)
- status (Text)
- name (Text, Required)
- email (Text, Required)
- phone (Text)
- linkedin_url (Text)
- experience (Text)
- skills (Text)
- summary (Text)
- resume_path (Text)
- metadata (JSONB)
```

#### `companies`
Companies discovered from Apollo API
```sql
- id (UUID, Primary Key)
- created_at (Timestamp)
- status (Text)
- apollo_id (Text, Indexed)
- name (Text)
- website (Text)
- industry (Text, Default: 'TECH')
- sub_industry (Text)
- tech_details (Text)
- company_details (Text)
- metadata (JSONB)
```

#### `contacts`
Contacts associated with companies
```sql
- id (UUID, Primary Key)
- created_at (Timestamp)
- company_id (UUID, Foreign Key)
- apollo_id (Text, Indexed)
- status (Text)
- name (Text)
- email_id (Text)
- phone_number (Text)
- linkedin_url (Text)
- role (Text)
- metadata (JSONB)
```

#### `email_logs`
Email generation and sending activity logs
```sql
- id (UUID, Primary Key)
- contact_id (UUID, Foreign Key)
- company_id (UUID, Foreign Key)
- status (Text) -- GENERATED, SENT, ERROR
- email_stage (Text) -- INITIAL, FOLLOWUP
- email_subject (Text)
- email_body (Text)
- attachment_details (JSONB)
- error_message (Text)
- metadata (JSONB)
- created_at (Timestamp)
```

#### `gmail_tokens`
Gmail OAuth tokens for email sending
```sql
- id (UUID, Primary Key)
- created_at (Timestamp)
- access_token (Text)
- refresh_token (Text)
- token_type (Text)
- expiry (Timestamp)
```

## 🔧 Configuration

### Target Settings
- **Countries**: Germany, Netherlands, Canada, Sweden, Finland, Norway, Ireland, UK, Luxembourg, UAE, Singapore, Australia
- **Company Sizes**: 1-10, 11-50, 51-200, 201-500 employees
- **Industries**: Technology, Software, IT, Artificial Intelligence
- **Target Roles**: CTO, Talent Acquisition, Technical Recruiter, Product Manager

### AI Email Generation
- **Model**: GPT-3.5-turbo (configurable)
- **Temperature**: 0.7 (configurable)
- **Max Tokens**: 512 (configurable)
- **Prompt**: Professional cold outreach with IIT Kharagpur background, company research, and skill highlighting

## 📁 Project Structure

```
cold_emailer/
├── api/                    # HTTP handlers and API logic
│   ├── handlers.go        # Main API endpoint handlers
│   ├── helpers.go         # Business logic helpers
│   └── models.go          # Request/response models
├── apollo/                # Apollo API integration
│   ├── client.go          # Apollo API client
│   ├── apollocompanydetails.go
│   ├── apollocompanysearchstructs.go
│   ├── apollocontactsearchrequest.go
│   └── apolloenrichcontact.go
├── cmd/server/            # Application entry point
│   └── main.go           # Main server setup
├── constants/             # Application constants
│   ├── constants.go      # Environment variables and settings
│   ├── errors.go         # Error definitions
│   ├── helpers.go        # Utility functions
│   └── prompts.go        # AI prompts
├── db/                    # Database connection and management
│   └── db.go             # Database initialization and migrations
├── dbmodels/              # Database models (one per table)
│   ├── companies/        # Company model operations
│   ├── contacts/         # Contact model operations
│   ├── email_logs/       # Email log model operations
│   ├── gmailtokens/      # Gmail token model operations
│   └── profileinfo/      # Profile model operations
├── frontend/              # Web dashboard
│   ├── public/           # Static HTML files
│   │   └── index.html    # Main dashboard page
│   ├── src/              # JavaScript and CSS
│   │   └── app.js        # Frontend application logic
│   ├── Dockerfile        # Frontend container configuration
│   └── nginx.conf        # Nginx configuration
├── gmail/                 # Gmail API integration
│   └── client.go         # Gmail OAuth and sending logic
├── migrations/            # Database migration files
│   ├── 20250707_150303_profile_info.sql
│   ├── 20250708_161223_gmail_tokens.sql
│   ├── 20250708_192659_companies.sql
│   ├── 20250709_012642_contacts.sql
│   └── 20250709_014140_email_outreach.sql
├── openai/                # OpenAI API integration
│   └── client.go         # OpenAI client for email generation
├── storage/               # File storage management
│   ├── storage.go        # File upload and management
│   └── storage_test.go   # Storage tests
├── utils/                 # Utility functions
│   └── helpers.go        # Common helper functions
├── docker-compose.yml     # Docker services configuration
├── Dockerfile            # Backend application container
├── go.mod               # Go module dependencies
└── README.md            # This file
```

## 🔄 Usage Workflow

### Complete Setup Process

1. **Initial Setup**
   ```bash
   # Start services
   docker-compose up -d --build
   
   # Verify health
   curl http://localhost:8000/health
   ```

2. **Profile Configuration**
   ```bash
   # Upload profile and resume
   curl -X POST http://localhost:8000/api/profile \
     -F "name=Your Name" \
     -F "email=your.email@example.com" \
     -F "phone=+1234567890" \
     -F "linkedin_url=https://linkedin.com/in/yourprofile" \
     -F "experience=5 years of software development..." \
     -F "skills=Go, Python, React, AWS..." \
     -F "summary=Experienced software engineer..." \
     -F "resume=@/path/to/your/resume.pdf"
   ```

3. **Gmail Authentication**
   ```bash
   # Visit in browser to complete OAuth
   open http://localhost:8000/gmail-auth-initiate
   ```

4. **Access Dashboard**
   ```bash
   # Open the web dashboard
   open http://localhost:3000
   ```
   The dashboard provides:
   - **Email Logs Tab**: View all sent emails with pagination and filtering
   - **Companies & Contacts Tab**: Browse companies and their associated contacts
   - **System Config Display**: Current targeting parameters
5. **Database Enrichment**
   ```bash
   # Enrich with companies and contacts
   curl -X POST http://localhost:8000/api/enrich-database \
     -H "Content-Type: application/json" \
     -d '{"count_new_companies": 50, "max_contacts_per_company": 5}'
   
   # Backfill company details
   curl -X POST http://localhost:8000/api/backfill-company-details
   ```

6. **Start Email Campaign**
   ```bash
   # Send initial outreach emails
   curl -X POST http://localhost:8000/api/send-few-initial-emails
   
   # Check logs via API or dashboard
   curl http://localhost:8000/api/logs
   # OR visit http://localhost:3000 and check Email Logs tab
   ```

7. **Follow-up Campaign**
   ```bash
   # Send follow-up emails (after some time)
   curl -X POST http://localhost:8000/api/send-few-follow-up-emails
   ```

### Monitoring and Management

```bash
# Check system status
curl http://localhost:8000/api/status

# View email logs via API
curl http://localhost:8000/api/email-logs?page=1&limit=20

# View companies via API
curl http://localhost:8000/api/companies

# Check database (direct SQL)
docker exec -it cold_emailer_db_1 psql -U coldemailer -d coldemailer
```

**Or use the web dashboard at http://localhost:3000 for a visual interface to:**
- Browse email logs with pagination and filtering
- Search and filter companies by name, industry, or website
- View detailed contact information organized by company
- Monitor system configuration and targeting parameters

## 🔒 Security & Best Practices

### API Key Management
- Store all API keys in `.env` file (never commit to git)
- Use environment variables in production
- Rotate keys regularly

### Database Security
- Use strong passwords for database
- Implement proper connection pooling
- Regular backups recommended

### OAuth Security
- Secure redirect URI configuration
- Token refresh handling
- Proper scope management for Gmail API

## 🚨 Rate Limiting & Quotas

### Apollo API
- Respects Apollo's rate limits
- Implements exponential backoff
- Batches requests efficiently

### OpenAI API
- Configurable temperature and max tokens
- Error handling for quota limits
- Retry logic for transient failures

### Gmail API
- OAuth token refresh handling
- Respects Gmail sending limits
- Proper attachment handling

## 🐛 Troubleshooting

### Common Issues

#### Database Connection Issues
```bash
# Check database status
docker-compose ps

# View database logs
docker-compose logs db

# Connect to database manually
docker exec -it cold_emailer_db_1 psql -U coldemailer -d coldemailer
```

#### Frontend Issues
```bash
# Check frontend container logs
docker-compose logs frontend

# Verify frontend is accessible
curl http://localhost:3000

# Check if API proxy is working
curl http://localhost:3000/api/companies
```

**Contacts not showing for a company?**

If you click a company card in the dashboard and see "No contacts found for this company":
- Ensure you have run the enrichment process (see step 5 in Usage Workflow) to fetch contacts for companies.
- If enrichment has been run and you still see no contacts, check backend logs for errors.
- You can also open the browser console (F12) to see debug logs for the contacts API response.

#### Gmail Authentication Issues
- Verify redirect URI matches exactly in Google Cloud Console
- Check OAuth scopes are properly configured
- Ensure Gmail API is enabled in Google Cloud Console

#### API Key Issues
- Verify all API keys are correctly set in `.env`
- Check API key permissions and quotas
- Test API keys independently

#### Migration Issues
```bash
# Generate new migration
go run cmd/server/main.go --generate-migration your_migration_name

# Check migration status in database
SELECT * FROM schema_migrations;
```

## 📊 Monitoring & Logging

### Log Levels
- **Error**: API failures, database errors
- **Info**: Successful operations, status updates
- **Debug**: Detailed operation flow

### Key Metrics to Monitor
- Email generation success rate
- Email sending success rate
- API response times
- Database connection health
- File storage usage
- Frontend dashboard usage

## 🔮 Future Enhancements

- [x] Web UI for campaign management
- [ ] Advanced email templates
- [ ] A/B testing for email content
- [ ] Analytics dashboard with charts
- [ ] Webhook support for status updates
- [ ] Multi-user support
- [ ] Email scheduling
- [ ] CRM integrations
- [ ] Real-time notifications
- [ ] Advanced filtering and search
- [ ] Export functionality

## 📄 License

MIT License - see LICENSE file for details

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📞 Support

For issues and questions:
- Create an issue in the GitHub repository
- Check the troubleshooting section above
- Review logs for detailed error information

---

**Note**: This system is designed for legitimate job application outreach. Please ensure compliance with applicable laws and email service provider terms of service when using this tool.

## 🔒 Security & Best Practices

### API Key Management
- Store all API keys in `.env` file (never commit to git)
- Use environment variables in production
- Rotate keys regularly

### Database Security
- Use strong passwords for database
- Implement proper connection pooling
- Regular backups recommended

### OAuth Security
- Secure redirect URI configuration
- Token refresh handling
- Proper scope management for Gmail API

## 🚨 Rate Limiting & Quotas

### Apollo API
- Respects Apollo's rate limits
- Implements exponential backoff
- Batches requests efficiently

### OpenAI API
- Configurable temperature and max tokens
- Error handling for quota limits
- Retry logic for transient failures

### Gmail API
- OAuth token refresh handling
- Respects Gmail sending limits
- Proper attachment handling

## 🐛 Troubleshooting

### Common Issues

#### Database Connection Issues
```bash
# Check database status
docker-compose ps

# View database logs
docker-compose logs db

# Connect to database manually
docker exec -it cold_emailer_db_1 psql -U coldemailer -d coldemailer
```

#### Gmail Authentication Issues
- Verify redirect URI matches exactly in Google Cloud Console
- Check OAuth scopes are properly configured
- Ensure Gmail API is enabled in Google Cloud Console

#### API Key Issues
- Verify all API keys are correctly set in `.env`
- Check API key permissions and quotas
- Test API keys independently

#### Migration Issues
```bash
# Generate new migration
go run cmd/server/main.go --generate-migration your_migration_name

# Check migration status in database
SELECT * FROM schema_migrations;
```

## 📊 Monitoring & Logging

### Log Levels
- **Error**: API failures, database errors
- **Info**: Successful operations, status updates
- **Debug**: Detailed operation flow

### Key Metrics to Monitor
- Email generation success rate
- Email sending success rate
- API response times
- Database connection health
- File storage usage

## 📄 License

MIT License - see LICENSE file for details

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📞 Support

For issues and questions:
- Create an issue in the GitHub repository
- Check the troubleshooting section above
- Review logs for detailed error information

---

**Note**: This system is designed for legitimate job application outreach. Please ensure compliance with applicable laws and email service provider terms of service when using this tool.


-- SQL migration
CREATE TABLE if not exists contacts (
    id uuid primary key, 
    created_at timestamp not null default now(),
    company_id uuid references companies(id),
    apollo_id text, 
    status text, 
    name text,
    email_id text, 
    phone_number text, 
    linkedin_url text, 
    role text, 
    metadata jsonb
);

CREATE INDEX IF NOT EXISTS idx_contacts_apollo_id ON contacts(apollo_id);
CREATE INDEX IF NOT EXISTS idx_contacts_name_email_id ON contacts(name, email_id);
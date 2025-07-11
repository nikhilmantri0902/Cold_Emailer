-- SQL migration
CREATE TABLE if not exists email_logs(
    id uuid primary key,
    contact_id uuid references contacts(id),
    company_id uuid references companies(id),
    status TEXT, 
    email_stage TEXT,
    email_subject TEXT, 
    email_body TEXT,
    attachment_details jsonb,
    error_message TEXT,
    metadata jsonb
);
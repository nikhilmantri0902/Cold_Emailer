-- SQL migration
CREATE TABLE if not exists contacts (
    id uuid primary key, 
    created_at timestamp not null default now(),
    company_id uuid references companies(id),
    status text, 
    name text,
    email_id text, 
    phone_number text, 
    linkedin_url text, 
    role text, 
    metadata jsonb
);
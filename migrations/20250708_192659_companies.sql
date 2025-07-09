-- SQL migration
CREATE TABLE if not exists companies (
    id uuid primary key, 
    created_at timestamp not null default now(),
    status text, 
    name text, 
    website text,
    industry text not null DEFAULT 'TECH', 
    sub_industry text,
    tech_details text, 
    company_details text,
    metadata jsonb
);
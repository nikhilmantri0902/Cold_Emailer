-- SQL migration
CREATE TABLE if not exists jobs (
    id uuid primary key, 
    created_at timestamp not null default now(),  
    completed_at timestamp,
    status text, 
    name text,
    message text, 
    metadata jsonb
);
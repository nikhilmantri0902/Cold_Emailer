CREATE TABLE if not exists profile_info (
    id uuid primary key, 
    created_at timestamp not null default now(),  
    status text, 
    name text not null, 
    email text not null, 
    phone text, 
    linkedin_url text, 
    experience text, 
    skills text, 
    summary text
);

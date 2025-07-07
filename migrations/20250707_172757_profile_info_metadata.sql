-- SQL migration
Alter table profile_info add column if not exists metadata jsonb;
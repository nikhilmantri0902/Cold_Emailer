-- SQL migration
Alter table email_logs add column if not exists created_at timestamp not null default now();


update email_logs set created_at = now() - interval '7 days' where email_stage = 'GENERATED';
update email_logs set created_at = now() - interval '4 days' where email_stage = 'SENT';

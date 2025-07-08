-- SQL migration
CREATE TABLE IF NOT EXISTS gmail_tokens (
    id uuid PRIMARY KEY,
    email_id text NOT NULL,
    access_token text NOT NULL,
    refresh_token text NOT NULL,
    expiry timestamp NOT NULL,
    created_at timestamp NOT NULL DEFAULT now()
);

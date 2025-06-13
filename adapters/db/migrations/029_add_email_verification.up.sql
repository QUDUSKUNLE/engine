-- Add email verification columns to users table
ALTER TABLE users ADD COLUMN email_verified boolean DEFAULT false;
ALTER TABLE users ADD COLUMN email_verified_at timestamp with time zone;

-- Create email verification tokens table
CREATE TABLE email_verification_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email text NOT NULL REFERENCES users(email) ON DELETE CASCADE,
    token text NOT NULL,
    used boolean DEFAULT false,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_verification_tokens_email ON email_verification_tokens(email);
CREATE INDEX idx_email_verification_tokens_token ON email_verification_tokens(token);

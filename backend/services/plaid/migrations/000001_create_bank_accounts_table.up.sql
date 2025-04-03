CREATE TABLE IF NOT EXISTS bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    plaid_account_id TEXT NOT NULL,
    access_token TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    subtype TEXT,
    mask TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, plaid_account_id)
);

CREATE INDEX IF NOT EXISTS idx_bank_accounts_user_id ON bank_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_plaid_account_id ON bank_accounts(plaid_account_id); 
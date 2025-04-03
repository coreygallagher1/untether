ALTER TABLE IF EXISTS plaid_items
ADD COLUMN IF NOT EXISTS webhook_url VARCHAR(255),
ADD COLUMN IF NOT EXISTS last_webhook_timestamp TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_plaid_items_last_webhook ON plaid_items(last_webhook_timestamp); 
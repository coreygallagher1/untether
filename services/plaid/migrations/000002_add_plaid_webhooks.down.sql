DROP INDEX IF EXISTS idx_plaid_items_last_webhook;

ALTER TABLE IF EXISTS plaid_items
DROP COLUMN IF EXISTS webhook_url,
DROP COLUMN IF EXISTS last_webhook_timestamp; 
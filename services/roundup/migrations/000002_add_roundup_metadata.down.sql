DROP INDEX IF EXISTS idx_roundups_category;
DROP INDEX IF EXISTS idx_roundups_merchant;

ALTER TABLE IF EXISTS roundups
DROP COLUMN IF EXISTS category,
DROP COLUMN IF EXISTS merchant_name,
DROP COLUMN IF EXISTS description; 
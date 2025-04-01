ALTER TABLE IF EXISTS roundups
ADD COLUMN IF NOT EXISTS description TEXT,
ADD COLUMN IF NOT EXISTS merchant_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS category VARCHAR(100);

CREATE INDEX IF NOT EXISTS idx_roundups_merchant ON roundups(merchant_name);
CREATE INDEX IF NOT EXISTS idx_roundups_category ON roundups(category); 
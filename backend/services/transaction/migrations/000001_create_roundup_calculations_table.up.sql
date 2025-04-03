CREATE TABLE IF NOT EXISTS roundup_calculations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    amount DECIMAL(10,2) NOT NULL,
    rounding_rule VARCHAR(50) NOT NULL,
    custom_rounding_amount DECIMAL(10,2),
    rounded_amount DECIMAL(10,2) NOT NULL,
    roundup_amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_roundup_calculations_created_at ON roundup_calculations(created_at); 
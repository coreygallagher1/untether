CREATE TABLE IF NOT EXISTS roundups (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    transaction_amount DECIMAL(10,2) NOT NULL,
    roundup_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_roundups_user_id ON roundups(user_id);
CREATE INDEX idx_roundups_status ON roundups(status); 
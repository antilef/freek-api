CREATE TABLE credit_cards (
    id SERIAL PRIMARY KEY,
    household_id INTEGER NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    closing_day INTEGER NOT NULL CHECK (closing_day BETWEEN 1 AND 28),
    member_id INTEGER REFERENCES members(id) ON DELETE SET NULL,
    credit_limit BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_credit_cards_household_id ON credit_cards(household_id);
CREATE INDEX idx_credit_cards_member_id ON credit_cards(member_id);

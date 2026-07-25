CREATE TYPE payment_type AS ENUM ('cash', 'debit', 'credit');

CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    household_id INTEGER NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    amount_clp BIGINT NOT NULL CHECK (amount_clp > 0),
    description VARCHAR(500) NOT NULL,
    notes TEXT,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    payment_type payment_type NOT NULL,
    credit_card_id INTEGER REFERENCES credit_cards(id) ON DELETE SET NULL,
    member_id INTEGER REFERENCES members(id) ON DELETE SET NULL,
    purchase_date DATE NOT NULL DEFAULT CURRENT_DATE,
    installments INTEGER NOT NULL DEFAULT 1 CHECK (installments >= 1),
    interest_rate NUMERIC(5,4) DEFAULT 0 CHECK (interest_rate >= 0),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_expenses_household_id ON expenses(household_id);
CREATE INDEX idx_expenses_credit_card_id ON expenses(credit_card_id);
CREATE INDEX idx_expenses_member_id ON expenses(member_id);
CREATE INDEX idx_expenses_category_id ON expenses(category_id);
CREATE INDEX idx_expenses_purchase_date ON expenses(purchase_date);
CREATE INDEX idx_expenses_payment_type ON expenses(payment_type);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    household_id INTEGER REFERENCES households(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    icon VARCHAR(50),
    is_predefined BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(household_id, name)
);

CREATE INDEX idx_categories_household_id ON categories(household_id);

-- Predefined categories (household_id = NULL means global)
INSERT INTO categories (name, icon, is_predefined) VALUES
    ('Alimentación', 'utensils', true),
    ('Transporte', 'car', true),
    ('Salud', 'heart', true),
    ('Educación', 'book', true),
    ('Entretenimiento', 'film', true),
    ('Hogar', 'home', true),
    ('Ropa', 'shirt', true),
    ('Servicios', 'wifi', true),
    ('Otros', 'more-horizontal', true);

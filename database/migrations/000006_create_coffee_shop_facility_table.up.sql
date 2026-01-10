CREATE TABLE IF NOT EXISTS "coffee_shop_facility" (
    id SERIAL PRIMARY KEY,
    facility_id INT REFERENCES facilities(id) ON DELETE CASCADE,
    coffee_shop_id INT REFERENCES coffee_shops(id) ON DELETE CASCADE
);
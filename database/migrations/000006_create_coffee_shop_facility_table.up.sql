CREATE TABLE IF NOT EXISTS "coffee_shop_facility" (
    id SERIAL PRIMARY KEY,
    facility_id INT REFERENCES facilities(id) NOT NULL,
    coffee_shop_id INT REFERENCES coffee_shops(id) NOT NULL
)
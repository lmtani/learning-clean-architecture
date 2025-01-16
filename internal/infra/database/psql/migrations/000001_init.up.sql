CREATE TABLE orders (
    id TEXT PRIMARY KEY,
    price NUMERIC NOT NULL,
    tax NUMERIC NOT NULL,
    final_price NUMERIC NOT NULL
);

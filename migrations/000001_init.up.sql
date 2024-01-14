CREATE TABLE warehouse (
    id SERIAL PRIMARY KEY,
    name varchar(150),
    is_available BOOLEAN
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name varchar(150),
    size varchar(150),
    code varchar(150) UNIQUE,
    quantity INTEGER NOT NULL,
    warehouse_id INTEGER REFERENCES warehouse(id)
);
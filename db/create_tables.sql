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

INSERT INTO warehouse (id, name, is_available) VALUES (1, 'Main Warehouse', true);

INSERT INTO warehouse (id, name, is_available) VALUES (2, 'Secondary Warehouse', true);



INSERT INTO products (id,name, size, code, quantity, warehouse_id) VALUES (1,'Футболка', '22', 'ABC123', 10, 1);

INSERT INTO products (id,name, size, code, quantity, warehouse_id) VALUES (2,'Джинсы', '44', 'DEF456', 5, 1);

INSERT INTO products (id,name, size, code, quantity, warehouse_id) VALUES (3,'Кроссовки', '35', 'GHI789', 8, 1);

INSERT INTO products (id,name, size, code, quantity, warehouse_id) VALUES (4,'Водолазка', '54', '123456', 50, 2);

INSERT INTO products (id,name, size, code, quantity, warehouse_id) VALUES (5,'Шлёпки', '66', '789012', 30, 2);

INSERT INTO products (id,name, size, code, quantity, warehouse_id) VALUES (6,'Брюки', '43', '345678', 20, 2);
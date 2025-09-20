CREATE TABLE IF NOT EXISTS ingredients (
    id INT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS inventory_items (
    ingredient_id INT PRIMARY KEY REFERENCES ingredients(id) ON DELETE RESTRICT ,
    quantity NUMERIC NOT NULL DEFAULT 0,
    unit TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS menu_items (
    id INT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    PRICE NUMERIC NOT NULL
);

CREATE TABLE IF NOT EXISTS menu_items_ingredients(
    menu_item_id INT NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    ingredient_id INT NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    quantity NUMERIC NOT NULL,
    PRIMARY KEY (menu_item_id, ingredient_id)
);



CREATE TABLE IF NOT EXISTS orders (
    id INT PRIMARY KEY,
    customer_name TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE IF NOT EXISTS order_items (
    order_item_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id INT NOT NULL REFERENCES menu_items(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (order_item_id, menu_item_id)
);
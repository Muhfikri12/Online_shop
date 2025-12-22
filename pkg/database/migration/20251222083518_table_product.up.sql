CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL,
    sku VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    category_id INT,
    price INT NOT NULL,
    stock INT NOT NULL,
    status VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP NULL
);

-- unique & index
CREATE UNIQUE INDEX idx_products_uuid ON products(uuid);
CREATE UNIQUE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_category_id ON products(category_id);

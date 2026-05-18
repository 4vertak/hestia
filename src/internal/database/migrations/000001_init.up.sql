CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    street VARCHAR(255) NOT NULL
);

CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    birthday DATE,
    gender CHAR(1) CHECK (gender IN ('M', 'F', 'O')),
    registration_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    address_id UUID REFERENCES addresses(id)
);

CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    address_id UUID REFERENCES addresses(id),
    phone_number VARCHAR(20)
);

CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image BYTEA NOT NULL
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    available_stock INT NOT NULL DEFAULT 0 CHECK (available_stock >= 0),
    last_update_date TIMESTAMPTZ,
    supplier_id UUID REFERENCES suppliers(id),
    image_id UUID REFERENCES images(id)
);

CREATE INDEX idx_clients_name_surname ON clients (surname, name);
CREATE INDEX idx_products_category ON products (category);
CREATE INDEX idx_products_available_stock ON products (available_stock) WHERE available_stock > 0;
CREATE INDEX idx_suppliers_name ON suppliers (name);
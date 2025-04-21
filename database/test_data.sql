CREATE TABLE users(
    id UUID PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    user_role VARCHAR(50) NOT NULL
);

CREATE TABLE pvz(
    id UUID PRIMARY KEY,
    registration_date TIMESTAMP NOT NULL,
    city VARCHAR(50) NOT NULL
);

CREATE TABLE receptions(
    id UUID PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    pvz_id UUID REFERENCES pvz(id) ON DELETE CASCADE,
    reception_status VARCHAR(20) NOT NULL
);

CREATE TABLE products(
    id UUID PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    product_type VARCHAR(30) NOT NULL,
    reception_id UUID REFERENCES receptions(id) ON DELETE CASCADE
);
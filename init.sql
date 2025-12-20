-- ============================================================
-- SPRINT 4: INICIALIZACIÓN DE BASE DE DATOS E-COMMERCE
-- Descripción: Schema Relacional + Datos Seed Realistas
-- ============================================================

-- 1. LIMPIEZA (Opcional, para reiniciar limpio si lo corres manual)
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS cart_items; -- Futuro: Si implementamos items en SQL
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;

-- 2. SCHEMA DEFINITION

-- A. Usuarios (Auth Service)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) DEFAULT 'hash_secreto', -- Simulación
    role VARCHAR(20) DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- B. Productos (Product Service - Catálogo Maestro)
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    category VARCHAR(50) NOT NULL, -- Electronics, Clothing, Home, Toys
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    stock INT DEFAULT 100,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- C. Pagos (Payment Service - Source of Truth)
CREATE TABLE payments (
    trace_id VARCHAR(100) PRIMARY KEY, -- Enlace con Logs/Traza
    user_id INT, -- FK opcional a users
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    gateway VARCHAR(50) NOT NULL, -- Stripe, PayPal
    status VARCHAR(20) NOT NULL, -- SUCCESS, FAILED, PENDING
    risk_score INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. SEED DATA: PRODUCTOS "HÉROE" (Realistas y Famosos)

INSERT INTO products (sku, name, category, price, stock) VALUES
-- ELECTRONICS (High Ticket)
('APP-IP15PM', 'Apple iPhone 15 Pro Max (256GB) - Titanium', 'Electronics', 1199.00, 50),
('APP-MBP16', 'Apple MacBook Pro 16 M3 Max', 'Electronics', 3499.00, 25),
('SAM-S24U', 'Samsung Galaxy S24 Ultra AI', 'Electronics', 1299.99, 60),
('SON-PS5-SL', 'Sony PlayStation 5 Slim Edition', 'Electronics', 499.99, 200),
('MIC-XSX', 'Microsoft Xbox Series X 1TB', 'Electronics', 499.99, 150),
('NIN-SW-OLED', 'Nintendo Switch OLED Model - Red/Blue', 'Electronics', 349.99, 300),
('NV-RTX4090', 'NVIDIA GeForce RTX 4090 Founders Edition', 'Electronics', 1599.00, 10),
('AMD-RX7900', 'AMD Radeon RX 7900 XTX', 'Electronics', 999.00, 20),
('LOG-MX3S', 'Logitech MX Master 3S Performance Mouse', 'Electronics', 99.99, 500),
('SON-WH1000', 'Sony WH-1000XM5 Noise Cancelling Headphones', 'Electronics', 348.00, 120),
('APP-APP2', 'Apple AirPods Pro (2nd Gen) USB-C', 'Electronics', 249.00, 400),
('DELL-AW34', 'Alienware 34 Curved QD-OLED Gaming Monitor', 'Electronics', 899.00, 30),

-- CLOTHING (Streetwear & Casual)
('NIKE-AF1', 'Nike Air Force 1 ''07 White', 'Clothing', 110.00, 1000),
('NIKE-J1-CHI', 'Air Jordan 1 Retro High OG "Chicago"', 'Clothing', 180.00, 15),
('ADI-SAMBA', 'Adidas Samba OG Cloud White', 'Clothing', 100.00, 250),
('NOR-NUPTSE', 'The North Face 1996 Retro Nuptse Jacket', 'Clothing', 320.00, 80),
('LEV-501', 'Levi''s 501 Original Fit Jeans', 'Clothing', 79.50, 400),
('PAT-FLEECE', 'Patagonia Better Sweater Fleece Jacket', 'Clothing', 149.00, 120),
('UNI-TSHIRT', 'Uniqlo U Crew Neck T-Shirt', 'Clothing', 19.90, 600),
('ZARA-COAT', 'Zara Wool Blend Coat', 'Clothing', 129.00, 90),
('VANS-OLD', 'Vans Old Skool Classic', 'Clothing', 70.00, 350),
('RAL-POLO', 'Ralph Lauren Custom Slim Fit Polo', 'Clothing', 98.00, 200),

-- HOME (Tech & Comfort)
('DYS-V15', 'Dyson V15 Detect Cordless Vacuum', 'Home', 749.99, 45),
('NES-VERTUO', 'Nespresso Vertuo Plus Coffee and Espresso Machine', 'Home', 159.00, 100),
('KIT-MIXER', 'KitchenAid Artisan Series 5-Quart Stand Mixer', 'Home', 449.99, 60),
('PHI-HUE-SK', 'Philips Hue White and Color Ambiance Starter Kit', 'Home', 199.99, 80),
('HER-MILLER', 'Herman Miller Aeron Chair', 'Home', 1695.00, 15),
('IKEA-MARKUS', 'IKEA Markus Office Chair', 'Home', 179.00, 300),
('STA-TUMBLER', 'Stanley Quencher H2.0 FlowState Tumbler 40oz', 'Home', 45.00, 500),
('YET-COOLER', 'YETI Tundra 45 Hard Cooler', 'Home', 325.00, 40),

-- TOYS (Collectibles & Kids)
('LEG-FALCON', 'LEGO Star Wars Millennium Falcon UCS', 'Toys', 849.99, 10),
('LEG-RIVEN', 'LEGO Lord of the Rings: Rivendell', 'Toys', 499.99, 20),
('FUN-POP-MAN', 'Funko Pop! Marvel: Iron Man', 'Toys', 12.99, 600),
('BAR-DREAM', 'Barbie Dreamhouse 2024', 'Toys', 199.00, 80),
('HOT-WHEELS', 'Hot Wheels 50-Car Pack', 'Toys', 54.99, 150),
('NERF-ELITE', 'Nerf Elite 2.0 Commander RD-6', 'Toys', 14.99, 400),
('HAS-MONO', 'Monopoly Classic Edition', 'Toys', 24.99, 300)
ON CONFLICT (sku) DO NOTHING;

-- 4. GENERACIÓN PROCEDURAL (Relleno Inteligente)
-- Generamos 300 productos adicionales combinando marcas y tipos para dar volumen

-- Generador Electronics
INSERT INTO products (sku, name, category, price, stock)
SELECT 
    'GEN-ELEC-' || generate_series(1, 150),
    (ARRAY['Samsung', 'LG', 'Asus', 'Acer', 'HP', 'Lenovo', 'Razer', 'Corsair'])[floor(random() * 8 + 1)] || ' ' ||
    (ARRAY['Ultra', 'Gaming', 'Pro', 'Office', 'Slim', '4K', 'Curved', 'Wireless'])[floor(random() * 8 + 1)] || ' ' ||
    (ARRAY['Monitor 27"', 'Keyboard', 'Mouse', 'Headset', 'Laptop', 'Router', 'Webcam', 'SSD 1TB'])[floor(random() * 8 + 1)],
    'Electronics',
    (random() * 800 + 50)::decimal(10,2),
    (random() * 50)::int
ON CONFLICT DO NOTHING;

-- Generador Clothing
INSERT INTO products (sku, name, category, price, stock)
SELECT 
    'GEN-CLOTH-' || generate_series(1, 150),
    (ARRAY['H&M', 'Zara', 'Gap', 'Uniqlo', 'Puma', 'Reebok', 'Under Armour', 'Columbia'])[floor(random() * 8 + 1)] || ' ' ||
    (ARRAY['Cotton', 'Summer', 'Winter', 'Sport', 'Vintage', 'Classic', 'Casual', 'Formal'])[floor(random() * 8 + 1)] || ' ' ||
    (ARRAY['T-Shirt', 'Hoodie', 'Shorts', 'Pants', 'Socks', 'Jacket', 'Cap', 'Sweater'])[floor(random() * 8 + 1)],
    'Clothing',
    (random() * 150 + 10)::decimal(10,2),
    (random() * 200)::int
ON CONFLICT DO NOTHING;

-- 5. USUARIOS FAKE (Para que Auth tenga sentido)
INSERT INTO users (username, email, role)
SELECT 
    'user_' || num,
    'user_' || num || '@example.com',
    CASE WHEN random() < 0.05 THEN 'admin' ELSE 'customer' END
FROM generate_series(1, 1000) AS num;

-- Fin del script
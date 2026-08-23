-- ============================================
-- Branch, user, role and position
-- ============================================

CREATE TABLE branches (
    branch_id VARCHAR(6) PRIMARY KEY,
    branch_name VARCHAR(100) NOT NULL,
    address VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE roles (
    role_id VARCHAR(20) PRIMARY KEY,
    role_name VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE positions (
    position_id VARCHAR(20) PRIMARY KEY,
    position_name VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    full_name VARCHAR(30) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    branch_id VARCHAR(6) NOT NULL,
    role_id VARCHAR(20) NOT NULL,
    position_id VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_users_branch
        FOREIGN KEY (branch_id)
        REFERENCES branches(branch_id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_users_role
        FOREIGN KEY (role_id)
        REFERENCES roles(role_id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_users_position
        FOREIGN KEY (position_id)
        REFERENCES positions(position_id)
        ON DELETE RESTRICT
);


-- ============================================
-- Master data
-- ============================================

CREATE TABLE brands (
    brand_id VARCHAR(5) PRIMARY KEY,
    brand_name VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE units (
    unit_id VARCHAR(5) PRIMARY KEY,
    unit_name VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE skus (
    sku_barcode VARCHAR(30) PRIMARY KEY,
    sku_name VARCHAR(50) NOT NULL,
    brand_id VARCHAR(5) NOT NULL,
    unit_id VARCHAR(5) NOT NULL,
    shelf_life_days INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_skus_brand
        FOREIGN KEY (brand_id)
        REFERENCES brands(brand_id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_skus_unit
        FOREIGN KEY (unit_id)
        REFERENCES units(unit_id)
        ON DELETE RESTRICT
);

CREATE TABLE stocks (
    branch_id VARCHAR(6) NOT NULL,
    sku_barcode VARCHAR(30) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (branch_id, sku_barcode),

    CONSTRAINT fk_stocks_branch
        FOREIGN KEY (branch_id)
        REFERENCES branches(branch_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_stocks_sku
        FOREIGN KEY (sku_barcode)
        REFERENCES skus(sku_barcode)
        ON DELETE CASCADE
);
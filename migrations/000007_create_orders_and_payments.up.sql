CREATE TABLE orders (
    order_id UUID PRIMARY KEY,
    branch_id VARCHAR(6) NOT NULL,
    cashier_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL,

    subtotal NUMERIC(12,2) NOT NULL DEFAULT 0
        CHECK (subtotal >= 0),

    discount_amount NUMERIC(12,2) NOT NULL DEFAULT 0
        CHECK (discount_amount >= 0),

    total_amount NUMERIC(12,2) NOT NULL DEFAULT 0
        CHECK (total_amount >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_orders_branch
        FOREIGN KEY (branch_id)
        REFERENCES branches(branch_id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_orders_cashier
        FOREIGN KEY (cashier_id)
        REFERENCES users(user_id)
        ON DELETE RESTRICT
);

CREATE TABLE order_items (
    order_id UUID NOT NULL,
    sku_barcode VARCHAR(30) NOT NULL,

    quantity NUMERIC(10,3) NOT NULL
        CHECK (quantity > 0),

    unit_price NUMERIC(12,2) NOT NULL,

    subtotal NUMERIC(12,2) NOT NULL
        CHECK (subtotal >= 0),

    discount_amount NUMERIC(12,2) NOT NULL DEFAULT 0
        CHECK (discount_amount >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (order_id, sku_barcode),

    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id)
        REFERENCES orders(order_id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_order_items_sku
        FOREIGN KEY (sku_barcode)
        REFERENCES skus(sku_barcode)
        ON DELETE RESTRICT
);

CREATE TABLE payments (
    payment_id UUID PRIMARY KEY,
    order_id UUID NOT NULL,

    method VARCHAR(20) NOT NULL,

    amount NUMERIC(12,2) NOT NULL
        CHECK (amount >= 0),

    status VARCHAR(20) NOT NULL,

    transaction_ref VARCHAR(100),

    paid_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_payments_order
        FOREIGN KEY (order_id)
        REFERENCES orders(order_id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_orders_branch_id
    ON orders(branch_id);

CREATE INDEX idx_payments_order_id
    ON payments(order_id);

CREATE INDEX idx_orders_open_branch
ON orders(branch_id, created_at DESC)
WHERE status = 'OPEN';
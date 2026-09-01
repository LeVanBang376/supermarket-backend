CREATE SEQUENCE import_request_id_seq
    START WITH 1
    INCREMENT BY 1;

CREATE TABLE import_requests (
    request_id VARCHAR(7) PRIMARY KEY,

    branch_id VARCHAR(6) NOT NULL,
    created_by UUID NOT NULL,

    expected_delivery_at TIMESTAMPTZ NOT NULL,
    delivery_license_plate VARCHAR(20),

    status VARCHAR(20) NOT NULL,

    received_by UUID,
    complete_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_import_requests_branch
        FOREIGN KEY (branch_id)
        REFERENCES branches(branch_id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_import_requests_created_by
        FOREIGN KEY (created_by)
        REFERENCES users(user_id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_import_requests_received_by
        FOREIGN KEY (received_by)
        REFERENCES users(user_id)
        ON DELETE RESTRICT
);


CREATE TABLE import_request_totes (
    request_id VARCHAR(7) NOT NULL,
    tote_barcode VARCHAR(30) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (request_id, tote_barcode),

    CONSTRAINT fk_import_request_totes_request
        FOREIGN KEY (request_id)
        REFERENCES import_requests(request_id)
        ON DELETE RESTRICT
);


CREATE TABLE import_request_products (
    request_id VARCHAR(7) NOT NULL,
    sku_barcode VARCHAR(30) NOT NULL,

    quantity INT NOT NULL
        CHECK (quantity > 0),

    tote_barcode VARCHAR(30),

    loaded_quantity INT NOT NULL DEFAULT 0
        CHECK (loaded_quantity >= 0),

    received_quantity INT NOT NULL DEFAULT 0
        CHECK (received_quantity >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (request_id, sku_barcode),

    CONSTRAINT fk_import_request_products_request
        FOREIGN KEY (request_id)
        REFERENCES import_requests(request_id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_import_request_products_sku
        FOREIGN KEY (sku_barcode)
        REFERENCES skus(sku_barcode)
        ON DELETE RESTRICT,

    CONSTRAINT fk_import_request_products_tote
        FOREIGN KEY (request_id, tote_barcode)
        REFERENCES import_request_totes(request_id, tote_barcode)
        ON DELETE RESTRICT
);

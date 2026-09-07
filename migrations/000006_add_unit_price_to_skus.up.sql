ALTER TABLE skus
ADD COLUMN unit_price NUMERIC(12,2) NOT NULL DEFAULT 0
    CHECK (unit_price >= 0);
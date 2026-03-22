-- +goose Up
CREATE TABLE orders IF NOT EXISTS (
    order_uuid UUID PRIMARY KEY, 
    user_uuid UUID NOT NULL,
    part_uuids UUID[] NOT NULL DEFAULT '{}',
    total_price DECIMAL(10, 2) NOT NULL,
    transaction_uuid UUID,
    payment_method VARCHAR(50),
    'status' VARCHAR(50) NOT NULL,
);

-- +goose Down
DROP TABLE orders IF EXISTS;

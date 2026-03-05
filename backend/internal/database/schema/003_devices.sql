-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE devices (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    device_name text NOT NULL,
    mac_address text NOT NULL,
    can_create bool NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE devices;

-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE doors (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    door_name text NOT NULL,
    device_id uuid REFERENCES devices(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE doors;

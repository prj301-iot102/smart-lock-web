-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE roles (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    role_name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE roles;

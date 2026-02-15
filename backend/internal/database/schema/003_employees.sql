-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE employees (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL REFERENCES users(id),
    full_name text UNIQUE NOT NULL,
    birth date NOT NULL,
    nfc_tag_id text NOT NULL REFERENCES nfc_tags(id),
    department text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE employess;

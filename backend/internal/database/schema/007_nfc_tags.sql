-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_cron WITH SCHEMA extensions;

CREATE TABLE nfc_tags (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    uid text NOT NULL,
    employee_id uuid REFERENCES employees(id),
    is_active boolean NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE nfc_tags;

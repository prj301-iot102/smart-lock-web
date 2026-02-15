CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +goose Up
CREATE TABLE nfc_tags (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    employee_id uuid NOT NULL REFERENCES employees(id),
    nfc_tag_id text NOT NULL,
    is_active boolean NOT NULL,
    enrolled_by uuid NOT NULL REFERENCES employees(id),
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE nfc_tags;

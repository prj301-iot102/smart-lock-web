-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE nfc_tags (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    uid text NOT NULL,
    is_active boolean NOT NULL,
    enrolled_by uuid NOT NULL REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE nfc_tags;

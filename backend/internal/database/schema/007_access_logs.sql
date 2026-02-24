-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE status AS ENUM ('granted', 'denied');

CREATE TABLE access_logs (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    employee_id uuid NOT NULL REFERENCES employees(id),
    door_id uuid NOT NULL REFERENCES doors(id),
    nfc_tag_id uuid NOT NULL REFERENCES nfc_tags(id),
    status status NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE access_logs;
DROP TYPE status;

-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE action AS ENUM ('enroll', 'update', 'revoke', 'delete');
CREATE TYPE result AS ENUM ('accepted', 'rejected', 'existed');

CREATE TABLE enrollment_logs (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    employee_id uuid NOT NULL REFERENCES employees(id),
    nfc_tag_uid text NOT NULL,
    action status NOT NULL,
    result result NOT NULL,
    admin_id uuid NOT NULL REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE enrollment_logs;
DROP TYPE action;

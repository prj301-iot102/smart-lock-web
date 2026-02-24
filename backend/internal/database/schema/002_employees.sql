-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE employees (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    role_id uuid NOT NULL REFERENCES roles(id),
    full_name text UNIQUE NOT NULL,
    birth date NOT NULL,
    department text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE employees;

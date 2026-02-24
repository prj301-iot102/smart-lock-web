-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE door_permissons (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    door_id uuid NOT NULL REFERENCES doors(id),
    role_id uuid NOT NULL REFERENCES roles(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE(door_id, role_id)
);

-- +goose Down
DROP TABLE door_permissons;

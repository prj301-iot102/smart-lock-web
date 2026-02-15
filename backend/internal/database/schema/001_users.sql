CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +goose Up
CREATE TABLE users (
	id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	username text UNIQUE NOT NULL,
	password text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE users;

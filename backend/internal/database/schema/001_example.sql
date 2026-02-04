-- +goose Up
CREATE TABLE examples (
	id int PRIMARY KEY NOT NULL,
	description text
);

-- +goose Down
DROP TABLE example;

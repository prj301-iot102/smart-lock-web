-- name: CreateUser :one
INSERT INTO users(username, password)
VALUES (@username, @password)
RETURNING id;

-- name: GetAccountByUsername :one
SELECT id, password
FROM users
WHERE username = @username;

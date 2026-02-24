-- name: CreateUser :one
INSERT INTO users(username, password, employee_id)
VALUES (@username, @password, employee_id)
RETURNING id;

-- name: GetAccountByUsername :one
SELECT password
FROM users u
WHERE username = @username;

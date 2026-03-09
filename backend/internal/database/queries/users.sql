-- name: CreateUser :one
INSERT INTO users(username, password, employee_id)
VALUES (@username, @password, employee_id)
RETURNING id;

-- name: GetAccountByUsername :one
SELECT id, password
FROM users u
WHERE username = @username;

-- name: GetAccountById :one
SELECT u.id, u.username, u.created_at, e.full_name, r.role_name
FROM users u
JOIN employees e ON e.id = u.employee_id
JOIN roles r ON r.id = e.role_id
WHERE u.id = @id;

-- name: UpdatePassword :exec
UPDATE users
SET
    password = @password
WHERE id = @id;

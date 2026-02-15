-- name: CreateUser :one
INSERT INTO users(username, password, employee_id)
VALUES (@username, @password, employee_id)
RETURNING id;

-- name: GetAccountByUsername :one
SELECT password, e.full_name
FROM users u
JOIN employees e ON e.id = u.employee_id
WHERE username = @username;

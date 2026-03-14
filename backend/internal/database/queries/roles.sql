-- name: SearchRoleName :many
SELECT * FROM roles
WHERE role_name LIKE '%' || @role_name || '%';

-- name: ListRoles :many
SELECT * FROM roles;

-- name: CreateRole :one
INSERT INTO roles (role_name)
VALUES(@role_name)
RETURNING id;

-- name: CreateEmployee :one
INSERT INTO employees (full_name, department)
VALUES ($1, $2)
RETURNING id;


-- name: GetEmployeeById :one
SELECT
    e.id,
    e.full_name,
    e.birth,
    r.role_name,
    e.department,
    e.created_at,
    e.updated_at
FROM employees e
JOIN roles r ON r.id = e.role_id
WHERE e.id = $1;


-- name: ListEmployees :many
SELECT
    id,
    full_name,
    birth,
    department,
    created_at,
    updated_at
FROM employees
WHERE
    ($1::text IS NULL OR full_name ILIKE '%' || $1 || '%') AND
    ($2::date IS NULL OR birth >= $2) AND
    ($3::date IS NULL OR birth <= $3) AND
    ($4::text IS NULL OR department ILIKE '%' || $4 || '%') AND
    ($5::date IS NULL OR created_at >= $5) AND
    ($6::date IS NULL OR created_at <= $6) AND
    ($7::date IS NULL OR updated_at >= $7) AND
    ($8::date IS NULL OR updated_at <= $8)
LIMIT $9 OFFSET $10;


-- name: UpdateEmployee :exec
UPDATE employees
SET
    full_name = COALESCE($1, full_name),
    birth = COALESCE($2, birth),
    department = COALESCE($3, department),
    role_id = COALESCE($4, role_id),
    updated_at = now()
WHERE id = $5;


-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE id = $1;

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
    (sqlc.narg(search)::text IS NULL OR full_name ILIKE '%' || $1::text || '%') AND
    (sqlc.narg(birth_from)::date IS NULL OR birth >= $2::date) AND
    (sqlc.narg(birth_to)::date IS NULL OR birth <= $3::date) AND
    (sqlc.narg(department)::text IS NULL OR department ILIKE '%' || $4::text || '%') AND
    (sqlc.narg(created_from)::date IS NULL OR created_at >= $5::date) AND
    (sqlc.narg(created_to)::date IS NULL OR created_at <= $6::date) AND
    (sqlc.narg(updated_from)::date IS NULL OR updated_at >= $7::date) AND
    (sqlc.narg(updated_to)::date IS NULL OR updated_at <= $8::date)
ORDER BY created_at DESC
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

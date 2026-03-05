-- name: CreateEmployee :one
INSERT INTO employees (full_name, department)
VALUES(@full_name, @department)
RETURNING id;

-- name: GetEmployeeById :one
SELECT full_name, birth, department, created_at, updated_at
FROM employees
WHERE id = @id;

-- name: UpdateEmployee :exec
UPDATE employees
SET
	full_name = COALESCE(sqlc.narg('full_name'), full_name),
	birth = COALESCE(sqlc.narg('birth'), birth),
	department = COALESCE(sqlc.narg('department'), department),
	updated_at = now()
WHERE id = @id;

SELECT id, full_name, birth, department, created_at, updated_at
FROM employees
WHERE
    ($1::text IS NULL OR full_name ILIKE '%' || $1 || '%')
    AND ($2::text IS NULL OR department = $2)
    AND ($3::date IS NULL OR birth >= $3)
    AND ($4::date IS NULL OR birth <= $4)
    AND ($5::date IS NULL OR created_at::date >= $5)
    AND ($6::date IS NULL OR created_at::date <= $6)
ORDER BY created_at DESC
LIMIT $7 OFFSET $8;

DELETE FROM employees
WHERE id = $1;


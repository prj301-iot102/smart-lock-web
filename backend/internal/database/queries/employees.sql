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

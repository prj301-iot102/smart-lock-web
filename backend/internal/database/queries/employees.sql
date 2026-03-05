-- name: CreateEmployee :one
INSERT INTO employees (full_name, department)
VALUES(@full_name, @department)
RETURNING id;

-- name: GetEmployeeById :one
SELECT e.id, e.full_name, e.birth, r.role_name, e.department, e.created_at, e.updated_at
FROM employees e
JOIN roles r ON r.id = e.role_id
WHERE e.id = @id;

-- name: UpdateEmployee :exec
UPDATE employees
SET
	full_name = COALESCE(sqlc.narg('full_name'), full_name),
	birth = COALESCE(sqlc.narg('birth'), birth),
	department = COALESCE(sqlc.narg('department'), department),
	updated_at = now()
WHERE id = @id;

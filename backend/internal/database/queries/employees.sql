-- name: CreateEmployee :one
INSERT INTO employees (full_name, nfc_tag_id, department)
VALUES(@full_name, @nfc_tag_id, @department)
RETURNING id;

-- name: GetEmployeeById :one
SELECT full_name, birth, nfc_tag_id, department, created_at, updated_at
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

-- name: UpdateTag :exec
UPDATE employees
SET
    nfc_tag_id = COALESCE(sqlc.narg('nfc_tag_id'), nfc_tag_id)
WHERE id = @id;

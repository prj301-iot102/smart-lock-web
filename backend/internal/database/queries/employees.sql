-- name: CreateEmployee :one
INSERT INTO employees (user_id, full_name, nfc_tag_id, department)
VALUES(@user_id, @full_name, @nfc_tag_id, @department)
RETURNING id;

-- name: GetEmployeeById :one
SELECT user_id, full_name, birth, nfc_tag_id, department, employees.created_at, updated_at
FROM employees
JOIN users ON employees.user_id = users.id
WHERE employees.id = @id;

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

-- name: CreateNfcTag :one
INSERT INTO nfc_tags (employee_id, uid, is_active, enrolled_by)
VALUES(@employee_id, @uid, @is_active, @enrolled_by)
RETURNING id;

-- name: UpdateTagStatus :exec
UPDATE nfc_tags
SET
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE id = @id;

-- name: GetTagById :one
SELECT *
FROM nfc_tags
JOIN employees ON employees.id = nfc_tags.employee_id
JOIN users ON users.id = nfc_tags.enrolled_by
WHERE nfc_tags.id = @id;

-- name: FilterTags :many
SELECT nfc_tags.uid, employees.full_name, is_active, nfc_tags.created_at
FROM nfc_tags
JOIN employees ON employees.id = nfc_tags.employee_id
WHERE
(
    (sqlc.narg('is_active')::boolean IS NULL OR is_active = @is_active)
);

-- name: CreateNfcTag :one
INSERT INTO nfc_tags (employee_id, uid, is_active, enrolled_by)
VALUES(@employee_id, @uid, @is_active, @enrolled_by)
RETURNING id;

-- name: UpdateTagStatus :one
UPDATE nfc_tags
SET
    is_active = @is_active
WHERE id = @id AND is_active = true
RETURNING id;

-- name: GetTagById :one
SELECT nt.id, nt.uid, nt.is_active, e.full_name, u.username, nt.created_at, nt.updated_at
FROM nfc_tags nt
JOIN employees e ON e.id = nt.employee_id
JOIN users u ON u.id = nt.enrolled_by
WHERE nt.id = @id;

-- name: FilterTags :many
SELECT nfc_tags.uid, employees.full_name, is_active, nfc_tags.created_at
FROM nfc_tags
JOIN employees ON employees.id = nfc_tags.employee_id
WHERE
(
    (sqlc.narg('is_active')::boolean IS NULL OR is_active = @is_active)
);

-- name: CreateNfcTag :one
INSERT INTO nfc_tags (uid)
VALUES(@uid)
RETURNING id;

-- name: UpdateTagStatus :one
UPDATE nfc_tags
SET
    is_active = @is_active
WHERE id = @id AND is_active = true
RETURNING id;

-- name: GetTagById :one
SELECT nt.id, nt.uid, nt.is_active, nt.employee_id, e.full_name, r.role_name, nt.created_at, nt.updated_at
FROM nfc_tags nt
JOIN employees e ON e.id = nt.employee_id
JOIN roles r ON r.id = e.role_id
WHERE nt.id = @id;

-- name: GetTagByUid :one
SELECT nt.id, nt.uid, e.full_name, nt.employee_id, nt.created_at, nt.updated_at
FROM nfc_tags nt
JOIN employees e ON e.id = nt.employee_id
WHERE nt.uid = @uid;

-- name: CheckUidExist :one
SELECT nt.id, nt.uid
FROM nfc_tags nt
WHERE nt.uid = @uid;

-- name: FilterTags :many
SELECT nfc_tags.uid, employees.full_name, is_active, nfc_tags.created_at
FROM nfc_tags
JOIN employees ON employees.id = nfc_tags.employee_id
WHERE
(
    (sqlc.narg('is_active')::boolean IS NULL OR is_active = @is_active)
);

-- name: CreateAccessLog :one
INSERT INTO access_logs (employee_id, status)
VALUES(@employee_id, @status)
RETURNING id;

-- name: GetAccessLogs :many
SELECT al.id, e.full_name, nt.uid, al.status, al.created_at
FROM access_logs al
JOIN employees e ON e.id = al.employee_id
JOIN nfc_tags nt ON nt.id = al.nfc_tag_id
ORDER BY al.created_at DESC;

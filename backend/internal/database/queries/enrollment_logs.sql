-- name: CreateEnrollmentLogs :one
INSERT INTO enrollment_logs (employee_id, nfc_tag_uid, action, admin_id)
VALUES (@employee_id, @nfc_tag_uid, @action, @admin_id)
RETURNING id;

-- name: GetEnrollmentLogs :many
SELECT el.id, e.full_name, el.nfc_tag_uid, el.action, el.created_at
FROM enrollment_logs el
JOIN employees e ON e.id = el.employee_id
JOIN users u ON u.id = el.employee_id
ORDER BY el.created_at DESC;

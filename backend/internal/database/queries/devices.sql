-- name: GetDeviceById :one
SELECT * FROM devices
WHERE id = @id;

-- name: UpdateDeviceCanCreate :one
UPDATE devices
SET
    can_create = @can_create
WHERE id = @id
RETURNING id;

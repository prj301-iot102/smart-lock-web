-- name: GetDeviceById :one
SELECT * FROM devices
WHERE id = @id;

-- name: GetDeviceByMac :one
SELECT * FROM devices
WHERE mac_address = @mac_address;

-- name: UpdateDeviceCanCreate :one
UPDATE devices
SET
    can_create = @can_create
WHERE id = @id
RETURNING id;

-- name: ListDevices :many
SELECT * FROM devices
ORDER BY device_name;

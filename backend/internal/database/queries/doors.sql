-- name: GetDoorById :one
SELECT * FROM doors
WHERE id = @id;

-- name: GetDoorByDeviceId :one
SELECT * FROM doors
WHERE device_id = @device_id;

-- name: GetDoorPermissonByDoorId :one
SELECT dp.id, dp.door_id, r.role_name, dp.created_at
FROM door_permissons dp
JOIN roles r ON r.id = dp.role_id
WHERE door_id = @door_id;

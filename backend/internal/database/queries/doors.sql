-- name: GetDoorById :one
SELECT * FROM doors d
WHERE id = @id;

-- name: GetDoorByDeviceId :one
SELECT * FROM doors
WHERE device_id = @device_id;

-- name: ListDoors :many
SELECT * FROM doors;

-- name: GetDoorPermissonByDoorId :many
SELECT r.role_name
FROM door_permissons dp
JOIN roles r ON r.id = dp.role_id
WHERE door_id = @door_id;

-- name: CheckDoorPermissons :one
SELECT dp.id
FROM door_permissons dp
JOIN roles r ON r.id = dp.role_id
WHERE door_id = @door_id AND role_id = @role_id;

-- name: AddDoorPermissionRole :one
INSERT INTO door_permissons(role_id, door_id)
VALUES(@role_id, @door_id)
RETURNING id;

-- name: DeleteDoorPermissonRole :exec
DELETE FROM door_permissons WHERE door_id = @door_id AND role_id = @role_id;


-- name: FindConnectionById :one
SELECT *
FROM connection
WHERE id = $1;

-- name: FindConnectionsByNodeId :many
SELECT *
FROM connection
WHERE node_id = $1;

-- name: SaveConnection :one
INSERT INTO connection (node_id, "user")
VALUES ($1, $2)
RETURNING id, created_at;

-- name: UpdateConnectionById :exec
UPDATE connection
SET "user" = $2
WHERE id = $1;

-- name: CheckExistsConnection :one
SELECT EXISTS (SELECT 1 FROM connection WHERE "user" = $1);

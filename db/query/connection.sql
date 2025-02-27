-- name: FindConnectionsById :many
SELECT *
FROM connection
WHERE node_id = $1;

-- name: SaveConnection :one
INSERT INTO connection (node_id, key, "user", checksum)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at;

-- name: UpdateConnectionById :exec
UPDATE connection
SET key    = $2,
    "user" = $3
WHERE id = $1;

-- name: CheckExistsConnection :one
SELECT EXISTS (SELECT 1 FROM connection WHERE checksum = $1);

-- name: FindNodes :many
SELECT *
FROM node
LIMIT $1 OFFSET $2;


-- name: FindNodeById :one
SELECT *
FROM node
WHERE id = $1;

-- name: SaveNode :exec
INSERT INTO node (host, port)
VALUES ($1, $2);

-- name: UpdateNodeById :exec
UPDATE node
SET host = $2,
    port = $3
WHERE id = $1;

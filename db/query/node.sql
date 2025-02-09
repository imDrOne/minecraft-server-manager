-- name: FindNodes :many
SELECT *
FROM node
LIMIT $1 OFFSET $2;


-- name: FindNodeById :one
SELECT *
FROM node
WHERE id = $1;

-- name: SaveNode :one
INSERT INTO node (host, port)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateNodeById :exec
UPDATE node
SET host = $2,
    port = $3
WHERE id = $1;


-- name: CountNode :one
SELECT count(*) FROM node;

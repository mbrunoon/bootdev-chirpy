-- name: CreateChirp :one
INSERT INTO chirps (body, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: AllChirps :many
SELECT * 
FROM chirps
ORDER BY created_at ASC;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;

-- name: FindChirp :one
SELECT *
FROM chirps
WHERE id = $1
LIMIT 1;
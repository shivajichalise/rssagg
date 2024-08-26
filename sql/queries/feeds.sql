-- name: CreateFeed :one
INSERT INTO feeds (id, user_id, name, url, created_at, updated_at) 
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

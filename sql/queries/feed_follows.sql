-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at) 
VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetFeedFollows :many
SELECT * FROM feed_follows WHERE user_id = $1;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows WHERE id = $1 AND user_id = $2;

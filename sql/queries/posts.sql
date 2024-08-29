-- name: CreatePost :one
INSERT INTO posts (id, feed_id, title, description, url, published_at, created_at, updated_at) 
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC 
LIMIT $2;

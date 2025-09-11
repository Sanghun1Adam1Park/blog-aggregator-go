-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollowsForUser :many
SELECT u.name AS user_name, f.name AS feed_name
FROM feed_follows AS ff
JOIN users AS u 
  ON ff.user_id = u.id
JOIN feeds AS f
  ON ff.feed_id = f.id
WHERE ff.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1
AND feed_id = $2;
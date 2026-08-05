-- name: FollowFeed :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id
LIMIT 1;

-- name: GetFeedFollowsByUserID :many
SELECT * FROM feed_follows 
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: UnfollowFeed :one
WITH deleted_follow AS (
    DELETE FROM feed_follows
    WHERE feed_follows.user_id = $2 
      AND feed_follows.feed_id = (SELECT id FROM feeds WHERE feeds.url = $1)
    RETURNING feed_id
)
SELECT feeds.*
FROM feeds
JOIN deleted_follow ON feeds.id = deleted_follow.feed_id;
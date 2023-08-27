-- notifications.sql

-- Create a notification
-- name: CreateNotification :one
INSERT INTO notifications (
  user_id,
  message,
  is_read
) VALUES (
  $1, $2, $3
) RETURNING *;

-- Get a notification by ID
-- name: GetNotificationByID :one
SELECT * FROM notifications
WHERE id = $1 LIMIT 1;

-- List notifications for a user with pagination
-- name: ListNotifications :many
SELECT * FROM notifications
WHERE user_id = $1
LIMIT $2 OFFSET $3;

-- Delete a notification by ID
-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = $1;

-- Mark a notification as read
-- name: MarkNotificationAsRead :exec
UPDATE notifications SET
  is_read = true
WHERE id = $1 RETURNING *;

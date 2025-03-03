-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT ? OFFSET ?;

-- name: CreateUser :one
INSERT INTO users (
  id, username, email, password_hash, full_name, role
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
  username = ?,
  email = ?,
  password_hash = ?,
  full_name = ?,
  role = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: GetJob :one
SELECT * FROM jobs
WHERE id = ? LIMIT 1;

-- name: ListJobs :many
SELECT * FROM jobs
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListJobsByOwner :many
SELECT * FROM jobs
WHERE owner_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateJob :one
INSERT INTO jobs (
  id, name, description, status, start_date, end_date, owner_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateJob :one
UPDATE jobs
SET
  name = ?,
  description = ?,
  status = ?,
  start_date = ?,
  end_date = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = ?;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = ? LIMIT 1;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListTasksByUser :many
SELECT * FROM tasks
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListTasksByJob :many
SELECT * FROM tasks
WHERE job_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListTasksByStatus :many
SELECT * FROM tasks
WHERE status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateTask :one
INSERT INTO tasks (
  id, title, description, status, priority, due_date, user_id, job_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET
  title = ?,
  description = ?,
  status = ?,
  priority = ?,
  due_date = ?,
  user_id = ?,
  job_id = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ?;

-- name: GetComment :one
SELECT * FROM comments
WHERE id = ? LIMIT 1;

-- name: ListCommentsByTask :many
SELECT * FROM comments
WHERE task_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateComment :one
INSERT INTO comments (
  id, content, user_id, task_id
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateComment :one
UPDATE comments
SET
  content = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = ?;

-- name: GetAttachment :one
SELECT * FROM attachments
WHERE id = ? LIMIT 1;

-- name: ListAttachmentsByTask :many
SELECT * FROM attachments
WHERE task_id = ?
ORDER BY created_at DESC;

-- name: ListAttachmentsByJob :many
SELECT * FROM attachments
WHERE job_id = ?
ORDER BY created_at DESC;

-- name: CreateAttachment :one
INSERT INTO attachments (
  id, filename, file_path, file_size, mime_type, task_id, job_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteAttachment :exec
DELETE FROM attachments
WHERE id = ?;

-- name: GetNotification :one
SELECT * FROM notifications
WHERE id = ? LIMIT 1;

-- name: ListNotificationsByUser :many
SELECT * FROM notifications
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListUnreadNotificationsByUser :many
SELECT * FROM notifications
WHERE user_id = ? AND is_read = 0
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateNotification :one
INSERT INTO notifications (
  id, type, content, user_id, reference_id, reference_type
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: MarkNotificationAsRead :one
UPDATE notifications
SET is_read = 1
WHERE id = ?
RETURNING *;

-- name: MarkAllNotificationsAsRead :exec
UPDATE notifications
SET is_read = 1
WHERE user_id = ? AND is_read = 0;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = ?;

-- name: CreateActivityLog :one
INSERT INTO activity_logs (
  id, action, entity_type, entity_id, details, user_id
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: ListActivityLogsByEntity :many
SELECT * FROM activity_logs
WHERE entity_type = ? AND entity_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListActivityLogsByUser :many
SELECT * FROM activity_logs
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

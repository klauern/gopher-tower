-- Database schema for Gopher Tower
-- AUTOMATICALLY GENERATED FROM MIGRATIONS
-- DO NOT EDIT DIRECTLY
--
-- Generated from migration files
--
-- This schema is used by sqlc to generate Go code
-- for database operations

CREATE TABLE tasks (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  status TEXT NOT NULL DEFAULT 'pending',
  priority TEXT DEFAULT 'medium',
  due_date TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id TEXT,
  job_id TEXT,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (job_id) REFERENCES jobs(id)
);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_job_id ON tasks(job_id);
CREATE TABLE jobs (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  status TEXT NOT NULL DEFAULT 'active',
  start_date TIMESTAMP,
  end_date TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  owner_id TEXT, command TEXT, arguments TEXT, stdout TEXT, stderr TEXT,
  FOREIGN KEY (owner_id) REFERENCES users(id)
);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_owner_id ON jobs(owner_id);
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  full_name TEXT,
  role TEXT DEFAULT 'user',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_login TIMESTAMP
);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE TABLE comments (
  id TEXT PRIMARY KEY,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id TEXT NOT NULL,
  task_id TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (task_id) REFERENCES tasks(id)
);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_task_id ON comments(task_id);
CREATE TABLE attachments (
  id TEXT PRIMARY KEY,
  filename TEXT NOT NULL,
  file_path TEXT NOT NULL,
  file_size INTEGER NOT NULL,
  mime_type TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  task_id TEXT,
  job_id TEXT,
  FOREIGN KEY (task_id) REFERENCES tasks(id),
  FOREIGN KEY (job_id) REFERENCES jobs(id)
);
CREATE INDEX idx_attachments_task_id ON attachments(task_id);
CREATE INDEX idx_attachments_job_id ON attachments(job_id);
CREATE TABLE notifications (
  id TEXT PRIMARY KEY,
  type TEXT NOT NULL,
  content TEXT NOT NULL,
  is_read BOOLEAN DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id TEXT NOT NULL,
  reference_id TEXT,
  reference_type TEXT,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_reference ON notifications(reference_id, reference_type);
CREATE TABLE activity_logs (
  id TEXT PRIMARY KEY,
  action TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id TEXT NOT NULL,
  details TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX idx_activity_logs_entity ON activity_logs(entity_type, entity_id);
CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
CREATE TABLE environments (
    -- Unique identifier for each environment
    id SERIAL PRIMARY KEY,
    -- Human-readable name of the environment (e.g., "production", "staging")
    name VARCHAR(255) NOT NULL,
    -- Timestamp when the environment was created
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Timestamp when the environment was last updated
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE env_vars (
    -- Unique identifier for each environment variable
    id SERIAL PRIMARY KEY,
    -- Reference to the environment this variable belongs to
    environment_id INTEGER NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    -- Name/key of the environment variable
    key VARCHAR(255) NOT NULL,
    -- Value of the environment variable (NULL if sensitive and stored in env_secrets)
    value TEXT,
    -- Flag indicating if this variable contains sensitive data
    is_sensitive BOOLEAN NOT NULL DEFAULT false,
    -- Timestamp when the variable was created
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Timestamp when the variable was last updated
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(environment_id, key)
);
CREATE TABLE env_secrets (
    -- Unique identifier for each secret
    id SERIAL PRIMARY KEY,
    -- Reference to the environment variable this secret belongs to
    env_var_id INTEGER NOT NULL REFERENCES env_vars(id) ON DELETE CASCADE,
    -- Encrypted value of the sensitive environment variable
    encrypted_value BYTEA NOT NULL,
    -- Initialization Vector used for encryption
    iv BYTEA NOT NULL,
    -- Timestamp when the secret was created
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Timestamp when the secret was last updated
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT one_secret_per_var UNIQUE(env_var_id)
);
CREATE INDEX idx_env_vars_environment_id ON env_vars(environment_id);
CREATE INDEX idx_env_secrets_env_var_id ON env_secrets(env_var_id);

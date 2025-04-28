-- Environments table stores different deployment or runtime environments (e.g., dev, staging, prod)
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

-- Environment variables table stores key-value pairs associated with environments
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

-- Secure storage for sensitive environment variable values
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

-- Index for faster lookups of environment variables by environment
CREATE INDEX idx_env_vars_environment_id ON env_vars(environment_id);
-- Index for faster lookups of secrets by environment variable
CREATE INDEX idx_env_secrets_env_var_id ON env_secrets(env_var_id);

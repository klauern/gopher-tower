CREATE TABLE environments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table for storing environment variables
CREATE TABLE env_vars (
    id SERIAL PRIMARY KEY,
    environment_id INTEGER NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value TEXT,
    is_sensitive BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(environment_id, key)
);

-- Table for storing encrypted sensitive values
CREATE TABLE env_secrets (
    id SERIAL PRIMARY KEY,
    env_var_id INTEGER NOT NULL REFERENCES env_vars(id) ON DELETE CASCADE,
    encrypted_value BYTEA NOT NULL,
    iv BYTEA NOT NULL,  -- Initialization Vector for encryption
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT one_secret_per_var UNIQUE(env_var_id)
);

-- Indexes for better query performance
CREATE INDEX idx_env_vars_environment_id ON env_vars(environment_id);
CREATE INDEX idx_env_secrets_env_var_id ON env_secrets(env_var_id);

-- Drop indexes first
DROP INDEX IF EXISTS idx_env_vars_environment_id;
DROP INDEX IF EXISTS idx_env_secrets_env_var_id;

-- Drop tables in reverse order of creation (respecting foreign key constraints)
DROP TABLE IF EXISTS env_secrets;
DROP TABLE IF EXISTS env_vars;
DROP TABLE IF EXISTS environments;

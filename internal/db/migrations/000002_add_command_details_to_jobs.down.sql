-- Remove command execution details from the jobs table
-- Note: Dropping columns might require specific SQLite versions or workarounds.
-- This assumes direct DROP COLUMN is supported.
ALTER TABLE jobs DROP COLUMN command;
ALTER TABLE jobs DROP COLUMN arguments;
ALTER TABLE jobs DROP COLUMN stdout;
ALTER TABLE jobs DROP COLUMN stderr;

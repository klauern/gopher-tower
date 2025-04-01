-- Add command execution details to the jobs table

-- The command to be executed (e.g., "python", "node", "go run")
ALTER TABLE jobs ADD COLUMN command TEXT;

-- JSON array or space-separated string of command arguments
ALTER TABLE jobs ADD COLUMN arguments TEXT;

-- Captured standard output from command execution
ALTER TABLE jobs ADD COLUMN stdout TEXT;

-- Captured standard error from command execution
ALTER TABLE jobs ADD COLUMN stderr TEXT;

-- Add command execution details to the jobs table
ALTER TABLE jobs ADD COLUMN command TEXT;
ALTER TABLE jobs ADD COLUMN arguments TEXT;
ALTER TABLE jobs ADD COLUMN stdout TEXT;
ALTER TABLE jobs ADD COLUMN stderr TEXT;

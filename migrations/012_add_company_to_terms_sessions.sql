-- Migration 012: add company column to terms_sessions
ALTER TABLE terms_sessions ADD COLUMN IF NOT EXISTS company VARCHAR(100) DEFAULT '';

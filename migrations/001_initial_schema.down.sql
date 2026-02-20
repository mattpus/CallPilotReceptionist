-- migrations/001_initial_schema.down.sql

-- Drop triggers
DROP TRIGGER IF EXISTS update_businesses_updated_at ON businesses;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_appointments_requested_date;
DROP INDEX IF EXISTS idx_appointments_status;
DROP INDEX IF EXISTS idx_appointments_business_id;
DROP INDEX IF EXISTS idx_appointments_call_id;

DROP INDEX IF EXISTS idx_transcripts_timestamp;
DROP INDEX IF EXISTS idx_transcripts_call_id;

DROP INDEX IF EXISTS idx_interactions_type;
DROP INDEX IF EXISTS idx_interactions_call_id;

DROP INDEX IF EXISTS idx_calls_started_at;
DROP INDEX IF EXISTS idx_calls_created_at;
DROP INDEX IF EXISTS idx_calls_status;
DROP INDEX IF EXISTS idx_calls_provider_call_id;
DROP INDEX IF EXISTS idx_calls_business_id;

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_business_id;

-- Drop tables
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS transcripts;
DROP TABLE IF EXISTS interactions;
DROP TABLE IF EXISTS calls;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS businesses;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";

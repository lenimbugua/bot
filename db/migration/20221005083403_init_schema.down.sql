ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS "sessions_response_id_fkey";
ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS "sessions_question_id_fkey";
ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS "sessions_profile_id_fkey";
ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS "sessions_chanell_id_fkey";
ALTER TABLE IF EXISTS "responses" DROP CONSTRAINT IF EXISTS "responses_question_id_fkey";
ALTER TABLE IF EXISTS "questions" DROP CONSTRAINT IF EXISTS "questions_parent_id_fkey";
ALTER TABLE IF EXISTS "questions" DROP CONSTRAINT IF EXISTS "questions_client_id_fkey";
ALTER TABLE IF EXISTS "profiles" DROP CONSTRAINT IF EXISTS "profiles_id_fkey";

DROP TABLE IF EXISTS "profile_responses";
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "responses";
DROP TABLE IF EXISTS "questions";
DROP TABLE IF EXISTS "profiles";
DROP TABLE IF EXISTS "clients";
DROP TABLE IF EXISTS "channels";

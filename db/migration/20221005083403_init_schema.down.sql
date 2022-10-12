ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS "sessions_response_id_fkey";
ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS "sessions_question_id_fkey";
ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS "sessions_user_id_fkey";
ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS "sessions_chanell_id_fkey";
ALTER TABLE IF EXISTS "responses" DROP CONSTRAINT IF EXISTS "responses_question_id_fkey";
ALTER TABLE IF EXISTS "questions" DROP CONSTRAINT IF EXISTS "questions_parent_id_fkey";
ALTER TABLE IF EXISTS "questions" DROP CONSTRAINT IF EXISTS "questions_client_id_fkey";
ALTER TABLE IF EXISTS "users" DROP CONSTRAINT IF EXISTS "users_id_fkey";

DROP TABLE IF EXISTS "user_responses";
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "responses";
DROP TABLE IF EXISTS "questions";
DROP TABLE IF EXISTS "user_company";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "companies";
DROP TABLE IF EXISTS "channels";

CREATE TABLE "channels" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "companies" (
  "id" bigserial PRIMARY KEY,
  "phone" varchar(15) UNIQUE NOT NULL,
  "name" varchar(25) NOT NULL,
  "email" varchar(25) UNIQUE NOT NULL
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "mobile" varchar UNIQUE NOT NULL,
  "password_hash" varchar NOT NULL,
  "name" varchar(20) NOT NULL,
  "company_id" bigint NOT NULL
);

CREATE TABLE "questions" (
  "id" bigserial PRIMARY KEY,
  "question" text NOT NULL,
  "company_id" bigint NOT NULL,
  "type" varchar NOT NULL,
  "parent_id" bigint NOT NULL,
  "channel_id" bigint NOT NULL,
  "next_question_id" bigint NOT NULL
);

CREATE TABLE "responses" (
  "id" bigserial PRIMARY KEY,
  "question_id" bigint NOT NULL,
  "response" text NOT NULL,
  "next_question_id" bigint NOT NULL
);

CREATE TABLE "sessions" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "chanell_id" bigint NOT NULL,
  "question_id" bigint NOT NULL,
  "response_id" bigint NOT NULL
);

CREATE TABLE "user_responses" (
  "id" bigserial PRIMARY KEY,
  "response_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "option_id" bigint NOT NULL,
  "question_id" bigint NOT NULL
);

ALTER TABLE "users" ADD FOREIGN KEY ("id") REFERENCES "companies" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "questions" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "questions" ADD FOREIGN KEY ("parent_id") REFERENCES "questions" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "responses" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("chanell_id") REFERENCES "channels" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("response_id") REFERENCES "responses" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

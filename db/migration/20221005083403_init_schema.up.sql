CREATE TABLE "channels" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "clients" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(25) NOT NULL
);

CREATE TABLE "profiles" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(20) NOT NULL,
  "client_id" bigint NOT NULL
);

CREATE TABLE "questions" (
  "id" bigserial PRIMARY KEY,
  "question" text NOT NULL,
  "client_id" bigint NOT NULL,
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
  "profile_id" bigint NOT NULL,
  "chanell_id" bigint NOT NULL,
  "question_id" bigint NOT NULL,
  "response_id" bigint NOT NULL
);

CREATE TABLE "profile_responses" (
  "id" bigserial PRIMARY KEY,
  "response_id" bigint NOT NULL,
  "profile_id" bigint NOT NULL,
  "option_id" bigint NOT NULL,
  "question_id" bigint NOT NULL
);

ALTER TABLE "profiles" ADD FOREIGN KEY ("id") REFERENCES "clients" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "questions" ADD FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "questions" ADD FOREIGN KEY ("parent_id") REFERENCES "questions" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "responses" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("chanell_id") REFERENCES "channels" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("response_id") REFERENCES "responses" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

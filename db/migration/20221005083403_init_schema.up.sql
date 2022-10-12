CREATE TABLE "channels" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "companies" (
  "id" bigserial PRIMARY KEY,
  "email" varchar(25) UNIQUE NOT NULL,
  "phone" varchar(15) UNIQUE NOT NULL,
  "name" varchar(25) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "phone" varchar UNIQUE NOT NULL,
  "password_hash" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00z',
  "name" varchar(20) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_company" (
  "user_id" bigint,
  "company_id" bigint,
  PRIMARY KEY ("company_id", "user_id")
);

CREATE TABLE "questions" (
  "id" bigserial PRIMARY KEY,
  "question" text NOT NULL,
  "company_id" bigint NOT NULL,
  "type" varchar NOT NULL,
  "parent_id" bigint NOT NULL,
  "channel_id" bigint NOT NULL,
  "next_question_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "responses" (
  "id" bigserial PRIMARY KEY,
  "question_id" bigint NOT NULL,
  "response" text NOT NULL,
  "next_question_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "channel_id" bigint NOT NULL DEFAULT 0,
  "question_id" bigint NOT NULL DEFAULT 0,
  "response_id" bigint NOT NULL DEFAULT 0,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_responses" (
  "id" bigserial PRIMARY KEY,
  "response_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "question_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX "q.company_id" ON "questions" ("company_id");

CREATE INDEX ON "questions" ("parent_id");

CREATE INDEX ON "questions" ("channel_id");

CREATE INDEX ON "responses" ("question_id");

CREATE INDEX ON "sessions" ("user_id");

CREATE INDEX ON "sessions" ("channel_id");

CREATE INDEX ON "sessions" ("question_id");

CREATE INDEX ON "sessions" ("response_id");

CREATE INDEX ON "user_responses" ("response_id");

CREATE INDEX ON "user_responses" ("user_id");

CREATE INDEX ON "user_responses" ("question_id");

ALTER TABLE "user_company" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "user_company" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "user_responses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "user_responses" ADD FOREIGN KEY ("response_id") REFERENCES "responses" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "user_responses" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "questions" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "questions" ADD FOREIGN KEY ("parent_id") REFERENCES "questions" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "responses" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("channel_id") REFERENCES "channels" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("response_id") REFERENCES "responses" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

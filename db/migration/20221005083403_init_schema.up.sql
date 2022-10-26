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
  "company_id" bigint NOT NULL,
  "password_hash" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00z',
  "name" varchar(20) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "bots" (
  "id" bigserial PRIMARY KEY,
  "title" text NOT NULL,
  "company_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "questions" (
  "id" bigserial PRIMARY KEY,
  "question" text NOT NULL,
  "bot_id" bigint NOT NULL,
  "type" varchar NOT NULL,
  "parent_id" bigint NOT NULL DEFAULT 0,
  "next_question_id" bigint NOT NULL DEFAULT 0,
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

CREATE INDEX ON "questions" ("bot_id");

CREATE INDEX ON "questions" ("parent_id");

CREATE INDEX ON "bots" ("company_id");

CREATE INDEX ON "responses" ("question_id");

CREATE INDEX ON "sessions" ("user_id");

CREATE INDEX ON "sessions" ("channel_id");

CREATE INDEX ON "sessions" ("question_id");

CREATE INDEX ON "sessions" ("response_id");

CREATE INDEX ON "user_responses" ("response_id");

CREATE INDEX ON "user_responses" ("user_id");

CREATE INDEX ON "user_responses" ("question_id");

ALTER TABLE "users" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "user_responses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "user_responses" ADD FOREIGN KEY ("response_id") REFERENCES "responses" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "user_responses" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE "questions" ADD FOREIGN KEY ("bot_id") REFERENCES "bots" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "questions" ADD FOREIGN KEY ("parent_id") REFERENCES "questions" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "bots" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "responses" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("channel_id") REFERENCES "channels" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("response_id") REFERENCES "responses" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;


-- initialize default channel
INSERT INTO channels (id, name)VALUES (0, 'DEFAULT_CHANNEL_NAME');

-- initialize default company
INSERT INTO companies  (id, email, phone, name)VALUES (0, 'DEFAULT@EMAIL.COM', '000000000000', 'DEFAULT_COMPANY_NAME');

-- initialize default bot
INSERT INTO bots  (id, title, company_id)VALUES (0, 'DEFAULT_BOT_TITLE', 0);

-- initialize default question
INSERT INTO questions  (id, question, bot_id, type, parent_id, next_question_id)VALUES (0, 'DEFAULT_QUESTION', 0, 'DEFAULT_QUESTION_TYPE', 0, 0);

-- initialize default response
INSERT INTO responses  (id, question_id, response, next_question_id)VALUES (0, 0, 'DEFAULT_RESPONSE', 0);

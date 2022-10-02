CREATE TABLE "users" (
  "id" bigserial unique,
  "username" varchar(255) NOT NULL UNIQUE,
  "password" varchar(255) NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
-- ALTER TABLE "accounts" ADD CONSTRAINT "user_id_currency_key" UNIQUE ("user_id", "currency");

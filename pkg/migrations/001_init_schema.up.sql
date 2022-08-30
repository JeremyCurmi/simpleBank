CREATE TABLE "accounts" (
  "id" bigserial unique,
  "name" varchar(255),
  "user_id" bigint,
  "balance" bigint,
  "currency" varchar(50),
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz DEFAULT (now()),
  constraint pk_name_user_id primary key (name, user_id)
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "sender_id" bigint,
  "receiver_id" bigint,
  "amount" bigint,
  "timestamp" timestamptz DEFAULT (now())
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint,
  "amount" bigint NOT NULL,
  "timestamp" timestamptz DEFAULT (now())
);

CREATE INDEX ON "accounts" ("user_id");

CREATE INDEX ON "transactions" ("account_id");

CREATE INDEX ON "transfers" ("sender_id");

CREATE INDEX ON "transfers" ("receiver_id");

CREATE INDEX ON "transfers" ("sender_id", "receiver_id");

COMMENT ON COLUMN "transactions"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "transactions" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("sender_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("receiver_id") REFERENCES "accounts" ("id");
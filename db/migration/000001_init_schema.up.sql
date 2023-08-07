

CREATE TABLE "users" (
 "id" bigserial PRIMARY KEY,
 "first_name" VARCHAR NOT NULL,
 "last_name" VARCHAR NOT NULL,
 "email" VARCHAR NOT NULL,
 "password_hash" VARCHAR NOT NULL,
 "created_at" timestamptz NOT NULL DEFAULT (now())

);

CREATE TABLE "trips" (
  "id" bigserial PRIMARY KEY,
  "title" VARCHAR NOT NULL,
  "user_id" bigint NOT NULL,
  "start_date" timestamptz NOT NULL,
  "end_date" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "fellow_travelers" (
  "id" bigserial PRIMARY KEY,
  "trip_id" bigint NOT NULL,
  "fellow_first_name" VARCHAR NOT NULL,
  "fellow_last_name" VARCHAR NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "trips" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "fellow_travelers" ADD FOREIGN KEY ("trip_id") REFERENCES "trips" ("id");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "trips" ("user_id");

CREATE INDEX ON "fellow_travelers" ("trip_id");


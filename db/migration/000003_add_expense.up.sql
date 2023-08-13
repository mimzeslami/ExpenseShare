CREATE TABLE "expenses" (
  "id" bigserial PRIMARY KEY,
  "trip_id" bigint NOT NULL,
  "payer_traveler_id" bigint NOT NULL,
  "amount" numeric NOT NULL,
  "description" VARCHAR NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "expenses" ADD FOREIGN KEY ("trip_id") REFERENCES "trips" ("id");
ALTER TABLE "expenses" ADD FOREIGN KEY ("payer_traveler_id") REFERENCES "fellow_travelers" ("id");

CREATE INDEX ON "expenses" ("trip_id");
CREATE INDEX ON "expenses" ("payer_traveler_id");
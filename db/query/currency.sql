-- currencies.sql

-- Create a currency
-- name: CreateCurrency :one
INSERT INTO currencies (
  code,
  name,
  symbol,
  exchange_rate,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- Get a currency by ID
-- name: GetCurrencyByID :one
SELECT * FROM currencies
WHERE id = $1 LIMIT 1;

-- List currencies with pagination
-- name: ListCurrencies :many
SELECT * FROM currencies
LIMIT $1 OFFSET $2;

-- Update a currency by ID
-- name: UpdateCurrency :one
UPDATE currencies SET
  code = $2,
  name = $3,
  symbol = $4,
  exchange_rate = $5,
  updated_at = $6
WHERE id = $1 RETURNING *;

-- Delete a currency by ID
-- name: DeleteCurrency :exec
DELETE FROM currencies
WHERE id = $1;

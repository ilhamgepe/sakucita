-- name: CreateTransaction :one
INSERT INTO transactions (
  id,
  donation_message_id,
  payer_user_id,
  payee_user_id,
  amount,
  fee_fixed,
  fee_percentage,
  fee_amount,
  net_amount,
  currency,
  status,
  external_reference
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transactions WHERE id = $1 LIMIT 1;

-- name: GetTransactionByDonationMessageID :one
SELECT * FROM transactions WHERE donation_message_id = $1 LIMIT 1;

-- name: GetTransactionByExternalReference :one
SELECT * FROM transactions WHERE external_reference = $1 LIMIT 1;
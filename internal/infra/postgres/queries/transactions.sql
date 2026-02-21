-- name: CreateTransaction :one
INSERT INTO transactions (
  id,
  donation_message_id,
  payment_channel_id,
  payer_user_id,
  payee_user_id,
  amount,
  gateway_fee_fixed,
  gateway_fee_percentage,
  gateway_fee_amount,
  platform_fee_fixed,
  platform_fee_percentage,
  platform_fee_amount,
  fee_fixed,
  fee_percentage,
  fee_amount,
  net_amount,
  currency,
  status,
  external_reference
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
) RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transactions WHERE id = $1 LIMIT 1;

-- name: GetTransactionByDonationMessageID :one
SELECT * FROM transactions WHERE donation_message_id = $1 LIMIT 1;

-- name: GetTransactionByExternalReference :one
SELECT * FROM transactions WHERE external_reference = $1 LIMIT 1;
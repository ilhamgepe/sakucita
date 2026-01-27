-- transaction status
CREATE TYPE transaction_status AS ENUM (
  'PENDING',
  'PAID',
  'FAILED',
  'EXPIRED',
  'REFUNDED'
);

CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT uuidv7(),

  -- actors
  payer_user_id UUID NULL,
  payee_user_id UUID NOT NULL,

  -- money (smallest unit)
  amount BIGINT NOT NULL CHECK (amount > 0),

  -- fee config (parameter snapshot)
  fee_fixed BIGINT NOT NULL DEFAULT 0,
  fee_percentage NUMERIC(10,6) NOT NULL DEFAULT 0,

  -- final result
  fee_amount BIGINT NOT NULL,
  net_amount BIGINT NOT NULL,

  -- rounding transparency
  rounding_amount NUMERIC(10,6) NOT NULL DEFAULT 0,
  rounding_strategy VARCHAR(30) NOT NULL DEFAULT 'FLOOR_TO_PLATFORM',

  currency CHAR(3) NOT NULL DEFAULT 'IDR',

  status transaction_status NOT NULL DEFAULT 'PENDING',

  -- external payment reference
  external_reference VARCHAR(100),

  meta JSONB NOT NULL DEFAULT '{}',

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  paid_at TIMESTAMPTZ,
  settled_at TIMESTAMPTZ,

  CHECK (fee_amount + net_amount = amount)
);

CREATE INDEX idx_transactions_payee
  ON transactions(payee_user_id);

CREATE INDEX idx_transactions_status
  ON transactions(status);

CREATE INDEX idx_transactions_external_reference
  ON transactions(external_reference);




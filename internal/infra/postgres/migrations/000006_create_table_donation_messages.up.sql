CREATE TYPE donation_media_type AS ENUM (
  'TEXT',
  'YOUTUBE',
  'GIF'
);

CREATE TABLE donation_messages (
  id UUID PRIMARY KEY,

  payee_user_id UUID NOT NULL,

  payer_user_id UUID,
  payer_name    VARCHAR(100) NOT NULL,
  email         VARCHAR(255) NOT NULL,
  message       TEXT NOT NULL,

  media_type donation_media_type NOT NULL,

  -- media (nullable)
  media_url TEXT,              -- youtube url / gif url
  media_start_seconds INT,


  max_play_seconds INT,         -- jumlah detik yang dapat di putar
  price_per_second BIGINT,      -- misal 500
  amount BIGINT NOT NULL,       -- charged_seconds * price_per_second
  currency CHAR(3) NOT NULL DEFAULT 'IDR',   -- IDR, dll

  status VARCHAR(20) NOT NULL DEFAULT 'CREATED',
  -- CREATED | PAID | PLAYED | CANCELED | REJECTED

  meta        JSONB,
  played_at   TIMESTAMPTZ,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

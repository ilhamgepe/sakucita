CREATE TYPE media_provider AS ENUM (
  'YOUTUBE'
);

CREATE TYPE donation_media_type AS ENUM (
  'NONE',
  'TTS',
  'YOUTUBE'
);

CREATE TABLE donation_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  transaction_id UUID NOT NULL
    REFERENCES transactions(id) ON DELETE CASCADE,

  streamer_user_id UUID NOT NULL,

  donor_name VARCHAR(100),
  message TEXT,
  is_anonymous BOOLEAN NOT NULL DEFAULT false,

  media_type donation_media_type NOT NULL DEFAULT 'NONE',

  -- TTS
  tts_language VARCHAR(10),
  tts_voice VARCHAR(50),

  -- Media (YouTube only for now)
  media_provider media_provider,
  media_video_id VARCHAR(50),
  media_start_seconds INT CHECK (media_start_seconds >= 0),
  media_end_seconds INT CHECK (media_end_seconds >= 0),
  media_duration_seconds INT CHECK (media_duration_seconds >= 0),

  -- FACT ONLY (no control state)
  played_at TIMESTAMPTZ,

  meta JSONB NOT NULL DEFAULT '{}',

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (transaction_id)
);

CREATE INDEX idx_donation_messages_streamer_unplayed
  ON donation_messages(streamer_user_id, created_at)
  WHERE played_at IS NULL;

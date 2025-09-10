-- =========================================
-- Extensions
-- =========================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;

-- =========================================
-- Enums
-- =========================================
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'auth_provider') THEN
    -- faqat email va google: soddalashtirilgan auth
    CREATE TYPE auth_provider AS ENUM ('email', 'google');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'match_status') THEN
    CREATE TYPE match_status AS ENUM ('queued','matched','canceled','completed','expired');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'cefr_level') THEN
    CREATE TYPE cefr_level AS ENUM ('A1','A2','B1','B2','C1','C2');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'message_type') THEN
    CREATE TYPE message_type AS ENUM ('text','voice','system');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'report_reason') THEN
    CREATE TYPE report_reason AS ENUM ('spam','abuse','nsfw','other');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'report_status') THEN
    CREATE TYPE report_status AS ENUM ('open','closed');
  END IF;
END $$;

-- =========================================
-- Users / Auth (email tasdiqlansa yoki google bo'lsa kifoya)
-- =========================================
CREATE TABLE users (
  id             uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  email          citext UNIQUE,                     -- email orqali kiruvchilar uchun
  phone          text UNIQUE,                       -- ixtiyoriy, hozir ishlatmasangiz ham mayli
  display_name   text NOT NULL,
  avatar_url     text,
  locale         text NOT NULL DEFAULT 'en',
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now(),
  deleted_at     timestamptz
);

-- Auth identifikatsiya: email (email_verified=true) yoki google (provider_uid bilan)
CREATE TABLE auth_identities (
  id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id         uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider        auth_provider NOT NULL,           -- 'email' yoki 'google'
  provider_uid    text NOT NULL,                    -- google sub yoki email address
  email_verified  boolean NOT NULL DEFAULT false,   -- email uchun tasdiq bayrog'i
  created_at      timestamptz NOT NULL DEFAULT now(),
  UNIQUE (provider, provider_uid)
);

-- Email OTP/verify tokenlar (ixtiyoriy, soddalashtirilgan)
CREATE TABLE auth_email_tokens (
  id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  email       citext NOT NULL,
  code        text NOT NULL,                        -- 6-8 raqamli OTP
  purpose     text NOT NULL,                        -- 'login' | 'verify'
  used        boolean NOT NULL DEFAULT false,
  expires_at  timestamptz NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX auth_email_tokens_email_idx ON auth_email_tokens (email, expires_at);

-- =========================================
-- Profile + 2FA (profil ichida ixtiyoriy 2 bosqichli himoya)
-- =========================================
CREATE TABLE profiles (
  user_id         uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  level           cefr_level,
  about           text,
  interests_txt   text,                 -- tezkor va soddalashtirilgan field
  timezone        text,
  last_seen_at    timestamptz,
  -- 2FA sozlamalari (ixtiyoriy): user o'zi yoqadi/o'chiradi
  twofa_enabled   boolean NOT NULL DEFAULT false,
  twofa_type      text,                 -- 'totp' (masalan) yoki NULL
  twofa_secret    text,                 -- TOTP secret (agar kerak bo'lsa)
  updated_at      timestamptz NOT NULL DEFAULT now()
);

-- Normalized interests (ixtiyoriy)
CREATE TABLE interests (
  id            serial PRIMARY KEY,
  slug          text UNIQUE NOT NULL,
  title         text NOT NULL
);

CREATE TABLE user_interests (
  user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  interest_id   int  NOT NULL REFERENCES interests(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, interest_id)
);

-- =========================================
-- Social controls
-- =========================================
CREATE TABLE blocks (
  blocker_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  blocked_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at    timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (blocker_id, blocked_id)
);

CREATE TABLE reports (
  id             uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reporter_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  reason         report_reason NOT NULL,
  note           text,
  status         report_status NOT NULL DEFAULT 'open',   -- 'open' / 'closed'
  handled_by     uuid REFERENCES users(id),
  handled_at     timestamptz,
  created_at     timestamptz NOT NULL DEFAULT now()
);

-- =========================================
-- Matchmaking & Sessions
-- =========================================
-- Queue real ishlashi Redis’da; bu jadval audit/log/statistika uchun.
CREATE TABLE match_attempts (
  id               uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id          uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  desired_level    cefr_level,
  desired_interest text,
  status           match_status NOT NULL,
  matched_with     uuid REFERENCES users(id),
  created_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX match_attempts_user_time_idx   ON match_attempts (user_id, created_at DESC);
CREATE INDEX match_attempts_status_time_idx ON match_attempts (status, created_at DESC);

CREATE TABLE sessions (
  id             uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  a_user_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  b_user_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  started_at     timestamptz NOT NULL DEFAULT now(),
  ended_at       timestamptz,
  duration_sec   int,
  topic          text,           -- ixtiyoriy mavzu
  a_rating       int,            -- 1..5 (optional; app layerda tekshirasiz)
  b_rating       int,
  a_notes        text,
  b_notes        text,
  created_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX sessions_a_time_idx ON sessions (a_user_id, started_at DESC);
CREATE INDEX sessions_b_time_idx ON sessions (b_user_id, started_at DESC);
CREATE INDEX sessions_time_idx   ON sessions (started_at);

-- =========================================
-- Messages (partitioned by month)
-- =========================================
CREATE TABLE messages (
  id            bigserial,
  session_id    uuid NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  sender_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type          message_type NOT NULL DEFAULT 'text',
  text          text,
  audio_url     text,        -- S3/MinIO URL
  created_at    timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Misol uchun: 2025-09 bo'lagini yaratib qo'yamiz
CREATE TABLE IF NOT EXISTS messages_2025_09 PARTITION OF messages
FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');

-- Tezkor indekslar
CREATE INDEX messages_session_time_idx ON messages USING btree (session_id, created_at DESC);
CREATE INDEX messages_sender_time_idx  ON messages USING btree (sender_id, created_at DESC);
CREATE INDEX messages_text_trgm_idx    ON messages USING gin (text gin_trgm_ops);

-- =========================================
-- Gamification / Streaks
-- =========================================
CREATE TABLE streaks (
  user_id        uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  current_streak int  NOT NULL DEFAULT 0,
  best_streak    int  NOT NULL DEFAULT 0,
  last_activity  date
);

CREATE TABLE badges (
  id            serial PRIMARY KEY,
  code          text UNIQUE NOT NULL,  -- e.g. FIRST_CHAT, SEVEN_DAYS
  title         text NOT NULL,
  description   text
);

CREATE TABLE user_badges (
  user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  badge_id      int  NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
  earned_at     timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, badge_id)
);

-- =========================================
-- Admin / Moderation / Audit
-- =========================================
CREATE TABLE audit_logs (
  id            bigserial PRIMARY KEY,
  actor_user_id uuid REFERENCES users(id),
  action        text NOT NULL,           -- e.g. 'BAN_USER', 'CLOSE_REPORT'
  target_id     uuid,
  meta          jsonb,
  created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX audit_logs_actor_time_idx ON audit_logs (actor_user_id, created_at DESC);
CREATE INDEX audit_logs_meta_gin_idx   ON audit_logs USING gin (meta);

-- =========================================
-- Foydali indekslar (qidiruv/filtrlash)
-- =========================================
CREATE INDEX IF NOT EXISTS users_email_idx        ON users (email);
CREATE INDEX IF NOT EXISTS users_phone_idx        ON users (phone);
CREATE INDEX IF NOT EXISTS profiles_level_idx     ON profiles (level);
CREATE INDEX IF NOT EXISTS profiles_last_seen_idx ON profiles (last_seen_at DESC);
CREATE INDEX IF NOT EXISTS reports_status_time_idx ON reports (status, created_at DESC);

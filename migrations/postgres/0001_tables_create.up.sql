-- Enable extensions
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto; 

-- USERS

CREATE TABLE IF NOT EXISTS users (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email          citext UNIQUE NOT NULL,
  password_hash  text,
  google_id      text UNIQUE,
  display_name   text NOT NULL CHECK (char_length(display_name) BETWEEN 1 AND 80),
  avatar_url     text,
  age            int CHECK (age IS NULL OR age BETWEEN 13 AND 120),
  gender         text CHECK (gender IN ('male','female') OR gender IS NULL),
  country_code   char(2) CHECK (country_code ~ '^[A-Z]{2}$' OR country_code IS NULL),
  native_lang    text,
  target_lang    text,
  level          smallint CHECK (level BETWEEN 1 AND 6 OR level IS NULL),
  about          text,
  timezone       text,
  email_verified boolean NOT NULL DEFAULT false,
  deleted_at     timestamptz,
  last_seen      timestamptz,
  role           text NOT NULL DEFAULT 'user' CHECK (role IN ('admin','user')),
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now()
);

-- EMAIL OTP TOKENS
CREATE TABLE IF NOT EXISTS auth_email_tokens (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email       citext NOT NULL,
  code        text   NOT NULL,
  purpose     text   NOT NULL CHECK (purpose IN ('login','verify','change_email')),
  used        boolean NOT NULL DEFAULT false,
  expires_at  timestamptz NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now()
);

-- INTERESTS
CREATE TABLE IF NOT EXISTS interests (
  id      serial PRIMARY KEY,
  slug    text UNIQUE NOT NULL,
  title   text NOT NULL
);

CREATE TABLE IF NOT EXISTS user_interests (
  user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  interest_id int  NOT NULL REFERENCES interests(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, interest_id)
);

-- USER SETTINGS
CREATE TABLE IF NOT EXISTS user_settings (
  user_id              uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  discoverable         boolean NOT NULL DEFAULT true,
  allow_messages       boolean NOT NULL DEFAULT true,
  notify_push          boolean NOT NULL DEFAULT true,
  notify_email         boolean NOT NULL DEFAULT false,
  created_at           timestamptz NOT NULL DEFAULT now(),
  updated_at           timestamptz
);

-- MATCH PREFERENCES
CREATE TABLE IF NOT EXISTS match_preferences (
  user_id          uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  target_lang      text,
  min_level        smallint CHECK (min_level BETWEEN 1 AND 6 OR min_level IS NULL),
  max_level        smallint CHECK (max_level BETWEEN 1 AND 6 OR max_level IS NULL),
  gender_filter    text CHECK (gender_filter IN ('male','female') OR gender_filter IS NULL),
  min_rating       smallint CHECK (min_rating BETWEEN 1 AND 5 OR min_rating IS NULL),
  countries_allow  text[],
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz,
  CHECK (min_level IS NULL OR max_level IS NULL OR min_level <= max_level)
);


-- FRIENDS / BLOCKS
CREATE TABLE IF NOT EXISTS friends (
  user_id         uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  friend_user_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at      timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, friend_user_id),
  CHECK (user_id <> friend_user_id)
);

CREATE TABLE IF NOT EXISTS blocks (
  blocker_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  blocked_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at  timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (blocker_id, blocked_id),
  CHECK (blocker_id <> blocked_id)
);

-- REPORTS
CREATE TABLE IF NOT EXISTS reports (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  reporter_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  reason         text NOT NULL,
  note           text,
  status         text NOT NULL DEFAULT 'open' CHECK (status IN ('open','reviewed','closed')),
  created_at     timestamptz NOT NULL DEFAULT now(),
  CHECK (reporter_id <> target_user_id)
);


-- AUTH SESSIONS
CREATE TABLE IF NOT EXISTS auth_sessions (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at    timestamptz NOT NULL DEFAULT now(),
  expires_at    timestamptz NOT NULL DEFAULT (now() + interval '30 days'),
  revoked_at    timestamptz,
  user_agent    text,
  ip_address    inet
);

-- DEVICE TOKENS
CREATE TABLE IF NOT EXISTS device_tokens (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  platform      text CHECK (platform IN ('android','ios','web') OR platform IS NULL),
  token         text NOT NULL,
  created_at    timestamptz NOT NULL DEFAULT now(),
  revoked_at    timestamptz,
  UNIQUE (user_id, token)
);


-- NOTIFICATIONS
CREATE TABLE IF NOT EXISTS notifications (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  kind          text NOT NULL,
  title         text,
  body          text,
  payload       jsonb,
  created_at    timestamptz NOT NULL DEFAULT now(),
  read_at       timestamptz
);


-- MATCH ATTEMPTS
CREATE TABLE IF NOT EXISTS match_attempts (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id          uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  desired_level    smallint,
  desired_language text,
  status           text NOT NULL CHECK (status IN ('queued','matched','canceled','completed','expired')),
  matched_with     uuid REFERENCES users(id),
  created_at       timestamptz NOT NULL DEFAULT now()
);

-- CALL SESSIONS (1-to-1)
CREATE TABLE IF NOT EXISTS sessions (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  a_user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  b_user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  started_at   timestamptz NOT NULL DEFAULT now(),
  ended_at     timestamptz,
  topic        text,
  state        text NOT NULL DEFAULT 'active' CHECK (state IN ('active','completed','canceled')),
  CHECK (a_user_id <> b_user_id),
  CHECK (ended_at IS NULL OR ended_at >= started_at)
);

-- SESSION FEEDBACK
CREATE TABLE IF NOT EXISTS session_feedback (
  session_id  uuid NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  rater_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  ratee_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  rating      int  NOT NULL CHECK (rating BETWEEN 1 AND 5),
  comment     text,
  created_at  timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (session_id, rater_id),
  CHECK (rater_id <> ratee_id)
);

-- MESSAGES (only session-based chat)
CREATE TABLE IF NOT EXISTS messages (
  id          bigserial PRIMARY KEY,
  session_id  uuid NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  sender_id   uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  kind        text NOT NULL, -- 'text'|'system' etc
  body        text,
  created_at  timestamptz NOT NULL DEFAULT now()
);

-- TRIGGERS
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger t JOIN pg_class c ON t.tgrelid = c.oid
    WHERE t.tgname = 'users_set_updated_at' AND c.relname = 'users'
  ) THEN
    CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE PROCEDURE trigger_set_updated_at();
  END IF;
END;
$$;


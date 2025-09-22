-- USERS (auth + profile)
CREATE TABLE IF NOT EXISTS users (
  id             uuid PRIMARY KEY,
  email          citext UNIQUE NOT NULL,                 
  password_hash  text,                                 
  google_id      text UNIQUE,                            
  display_name   text NOT NULL CHECK (length(display_name) BETWEEN 1 AND 80),
  avatar_url     text,
  age            int CHECK (age IS NULL OR age BETWEEN 13 AND 120),
  gender         text CHECK (gender IN ('male','female') OR gender IS NULL),
  country_code   char(2) CHECK (country_code ~ '^[A-Z]{2}$' OR country_code IS NULL),
  native_lang    text,                                  
  target_lang    text,                                   
  level          smallint CHECK (level BETWEEN 1 AND 6 OR level IS NULL), 
  about          text,
  timezone       text,
  created_at     timestamptz NOT NULL DEFAULT now()
);



-- EMAIL OTP / VERIFY
CREATE TABLE IF NOT EXISTS auth_email_tokens (
  id          uuid PRIMARY KEY,
  email       citext NOT NULL,
  code        text   NOT NULL,                          
  purpose     text   NOT NULL CHECK (purpose IN ('login','verify','change_email')),
  used        boolean NOT NULL,
  expires_at  timestamptz NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now()
);

-- INTERESTS (oddiy)
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

-- USER SETTINGS (privacy/notify)
CREATE TABLE IF NOT EXISTS user_settings (
  user_id              uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  discoverable         boolean NOT NULL DEFAULT true,    
  allow_messages       boolean NOT NULL DEFAULT true,    
  notify_push          boolean NOT NULL DEFAULT true,    
  notify_email         boolean NOT NULL DEFAULT false,   
  created_at           timestamptz NOT NULL DEFAULT now(),
  updated_at           timestamptz
);

-- MATCH PREFERENCES (filterlar)

CREATE TABLE IF NOT EXISTS match_preferences (
  user_id          uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  target_lang      text,                -- qaysi til bo'yicha sherik
  min_level        smallint CHECK (min_level BETWEEN 1 AND 6 OR min_level IS NULL),
  max_level        smallint CHECK (max_level BETWEEN 1 AND 6 OR max_level IS NULL),
  gender_filter    text CHECK (gender_filter IN ('male','female') OR gender_filter IS NULL),
  min_rating       smallint CHECK (min_rating BETWEEN 1 AND 5 OR min_rating IS NULL),
  countries_allow  text[],            
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz,
  CHECK (min_level IS NULL OR max_level IS NULL OR min_level <= max_level)
);

-- SOCIAL: FRIENDS (yo'nalishli)
CREATE TABLE IF NOT EXISTS friends (
  user_id         uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  friend_user_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at      timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, friend_user_id),
  CHECK (user_id <> friend_user_id)
);

-- SOCIAL: BLOCK / REPORT
CREATE TABLE IF NOT EXISTS blocks (
  blocker_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  blocked_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at  timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (blocker_id, blocked_id),
  CHECK (blocker_id <> blocked_id)
);

CREATE TABLE IF NOT EXISTS reports (
  id             uuid PRIMARY KEY,
  reporter_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  reason         text NOT NULL,                      
  note           text,
  status         text NOT NULL DEFAULT 'open',        
  created_at     timestamptz NOT NULL DEFAULT now(),
  CHECK (reporter_id <> target_user_id)
);

-- AUTH SESSIONS (login sessiyalar)
CREATE TABLE IF NOT EXISTS auth_sessions (
  id            uuid PRIMARY KEY,
  user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at    timestamptz NOT NULL DEFAULT now(),
  expires_at    timestamptz,
  revoked_at    timestamptz,
  user_agent    text,
  ip_address    inet
);

-- DEVICES (push tokenlar)
CREATE TABLE IF NOT EXISTS device_tokens (
  id            uuid PRIMARY KEY,
  user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  platform      text,                                   -- 'for   flutter'
  token         text NOT NULL,                          -- FCM/APNs/WebPush token
  created_at    timestamptz NOT NULL DEFAULT now(),
  revoked_at    timestamptz,
  UNIQUE (user_id, token)                                 
);

-- NOTIFICATIONS (server → user)
CREATE TABLE IF NOT EXISTS notifications (
  id            uuid PRIMARY KEY,
  user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  kind          text NOT NULL,                          -- 'match','message','call',...
  title         text,
  body          text,
  payload       jsonb,                                  
  created_at    timestamptz NOT NULL DEFAULT now(),
  read_at       timestamptz
);

-- MATCH/AUDIT (log; servis hisoblaydi)
CREATE TABLE IF NOT EXISTS match_attempts (
  id               uuid PRIMARY KEY,
  user_id          uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  desired_level    smallint,                           
  desired_language text,                                -
  status           text NOT NULL,                       -- 'queued','matched','canceled','completed','expired'
  matched_with     uuid REFERENCES users(id),
  created_at       timestamptz NOT NULL DEFAULT now()
);

-- CALL SESSIONS (1-ga-1; media saqlanmaydi)
CREATE TABLE IF NOT EXISTS sessions (
  id           uuid PRIMARY KEY,
  a_user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  b_user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  started_at   timestamptz NOT NULL DEFAULT now(),
  ended_at     timestamptz,                              
  topic        text,
  CHECK (a_user_id <> b_user_id),
  CHECK (ended_at IS NULL OR ended_at >= started_at)
);

-- SESSION FEEDBACK (har foydalanuvchidan bittadan)
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


-- DIALOGS (offline/online chat guruhlash)

CREATE TABLE IF NOT EXISTS dialogs (
  id            uuid PRIMARY KEY,
  kind          text NOT NULL,                           
  created_at    timestamptz NOT NULL DEFAULT now(),
  created_by    uuid REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS dialog_members (
  dialog_id     uuid NOT NULL REFERENCES dialogs(id) ON DELETE CASCADE,
  user_id       uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  joined_at     timestamptz NOT NULL DEFAULT now(),
  role          text,                                    -- 'owner'|'member' (servis)
  PRIMARY KEY (dialog_id, user_id)
);

-- MESSAGES (faqat matn; audio saqlanmaydi)
CREATE TABLE IF NOT EXISTS messages (
  id          bigserial PRIMARY KEY,
  session_id  uuid REFERENCES sessions(id) ON DELETE CASCADE,
  dialog_id   uuid REFERENCES dialogs(id)  ON DELETE CASCADE,
  sender_id   uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  kind        text NOT NULL,                              -- 'text'|'system' 
  body        text,
  created_at  timestamptz NOT NULL DEFAULT now(),
  CHECK (
    (session_id IS NOT NULL AND dialog_id IS NULL) OR
    (session_id IS NULL AND dialog_id IS NOT NULL)
  )
);

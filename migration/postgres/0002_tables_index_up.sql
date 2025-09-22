-- USERS

CREATE INDEX IF NOT EXISTS users_match_idx
  ON users (target_lang, level, country_code, gender);
CREATE INDEX IF NOT EXISTS users_created_idx
  ON users (created_at DESC);

-- AUTH EMAIL TOKENS (OTP)

CREATE INDEX IF NOT EXISTS auth_email_tokens_email_code_idx
  ON auth_email_tokens (email, code);
-- Aktiv tokenlarni (used = false) va muddati bo'yicha tez topish
CREATE INDEX IF NOT EXISTS auth_email_tokens_lookup_idx
  ON auth_email_tokens (email, used, expires_at DESC);

-- INTERESTS
CREATE INDEX IF NOT EXISTS user_interests_interest_idx
  ON user_interests (interest_id, user_id);

-- USER SETTINGS / MATCH PREFERENCES
CREATE INDEX IF NOT EXISTS match_preferences_target_idx
  ON match_preferences (target_lang, min_level, max_level);

-- FRIENDS / BLOCKS

CREATE INDEX IF NOT EXISTS friends_reverse_idx
  ON friends (friend_user_id, user_id);
CREATE INDEX IF NOT EXISTS blocks_blocked_idx
  ON blocks (blocked_id, blocker_id);

-- REPORTS (moderatsiya)
CREATE INDEX IF NOT EXISTS reports_target_time_idx
  ON reports (target_user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS reports_status_time_idx
  ON reports (status, created_at DESC);

-- AUTH SESSIONS

CREATE INDEX IF NOT EXISTS auth_sessions_user_time_idx
  ON auth_sessions (user_id, created_at DESC);
-- Cron/cleanup uchun
CREATE INDEX IF NOT EXISTS auth_sessions_expires_idx
  ON auth_sessions (expires_at);

-- DEVICE TOKENS
CREATE INDEX IF NOT EXISTS device_tokens_user_idx
  ON device_tokens (user_id);
-- Tokenni topib o'chirish/yangilash uchun
CREATE INDEX IF NOT EXISTS device_tokens_token_idx
  ON device_tokens (token);

-- NOTIFICATIONS

CREATE INDEX IF NOT EXISTS notifications_user_time_idx
  ON notifications (user_id, created_at DESC);
-- Faqat "o'qilmaganlar" ni tez olish (partial index)
CREATE INDEX IF NOT EXISTS notifications_unread_idx
  ON notifications (user_id, created_at DESC)
  WHERE read_at IS NULL;

-- MATCH ATTEMPTS (audit/log)

CREATE INDEX IF NOT EXISTS match_attempts_user_time_idx
  ON match_attempts (user_id, created_at DESC);
-- Status bo'yicha boshqaruv paneli
CREATE INDEX IF NOT EXISTS match_attempts_status_time_idx
  ON match_attempts (status, created_at DESC);
-- Kim bilan match bo'lganini ko'rish
CREATE INDEX IF NOT EXISTS match_attempts_matched_with_idx
  ON match_attempts (matched_with, created_at DESC);


-- SESSIONS (qo'ng'iroqlar)

CREATE INDEX IF NOT EXISTS sessions_a_time_idx
  ON sessions (a_user_id, started_at DESC);
CREATE INDEX IF NOT EXISTS sessions_b_time_idx
  ON sessions (b_user_id, started_at DESC);
-- So'nggi qo'ng'iroqlar umumiy ko'rinishi
CREATE INDEX IF NOT EXISTS sessions_started_desc_idx
  ON sessions (started_at DESC);

-- SESSION FEEDBACK

CREATE INDEX IF NOT EXISTS session_feedback_ratee_time_idx
  ON session_feedback (ratee_id, created_at DESC);
-- O'zi bergan fikrlarni ko'rish
CREATE INDEX IF NOT EXISTS session_feedback_rater_time_idx
  ON session_feedback (rater_id, created_at DESC);
-- Sessiya bo'yicha join'lar
CREATE INDEX IF NOT EXISTS session_feedback_session_idx
  ON session_feedback (session_id);

-- DIALOGS / DIALOG MEMBERS

CREATE INDEX IF NOT EXISTS dialog_members_user_idx
  ON dialog_members (user_id, dialog_id);
-- "Dialogdagi barcha a'zolar" uchun PK (dialog_id, user_id) yetadi

-- MESSAGES (chat)

CREATE INDEX IF NOT EXISTS messages_session_paging_idx
  ON messages (session_id, id DESC);
-- Dialog (offline/online) chatini paginatsiya
CREATE INDEX IF NOT EXISTS messages_dialog_paging_idx
  ON messages (dialog_id, id DESC);
-- (ixtiyoriy) foydalanuvchi yuborgan xabarlar oqimini ko'rish
CREATE INDEX IF NOT EXISTS messages_sender_time_idx
  ON messages (sender_id, created_at DESC);

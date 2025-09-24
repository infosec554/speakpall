-- USERS
DROP INDEX IF EXISTS users_match_idx;
DROP INDEX IF EXISTS users_created_idx;

-- AUTH EMAIL TOKENS
DROP INDEX IF EXISTS auth_email_tokens_email_code_idx;
DROP INDEX IF EXISTS auth_email_tokens_lookup_idx;

-- INTERESTS
DROP INDEX IF EXISTS user_interests_interest_idx;

-- MATCH PREFERENCES
DROP INDEX IF EXISTS match_preferences_target_idx;

-- FRIENDS / BLOCKS
DROP INDEX IF EXISTS friends_reverse_idx;
DROP INDEX IF EXISTS blocks_blocked_idx;

-- REPORTS
DROP INDEX IF EXISTS reports_target_time_idx;
DROP INDEX IF EXISTS reports_status_time_idx;

-- AUTH SESSIONS
DROP INDEX IF EXISTS auth_sessions_user_time_idx;

-- DEVICE TOKENS
DROP INDEX IF EXISTS device_tokens_user_idx;

-- NOTIFICATIONS
DROP INDEX IF EXISTS notifications_user_time_idx;
DROP INDEX IF EXISTS notifications_unread_idx;

-- MATCH ATTEMPTS
DROP INDEX IF EXISTS match_attempts_user_time_idx;
DROP INDEX IF EXISTS match_attempts_status_time_idx;
DROP INDEX IF EXISTS match_attempts_matched_with_idx;

-- SESSIONS
DROP INDEX IF EXISTS sessions_a_time_idx;
DROP INDEX IF EXISTS sessions_b_time_idx;
DROP INDEX IF EXISTS sessions_started_desc_idx;

-- SESSION FEEDBACK
DROP INDEX IF EXISTS session_feedback_ratee_time_idx;
DROP INDEX IF EXISTS session_feedback_rater_time_idx;
DROP INDEX IF EXISTS session_feedback_session_idx;

-- MESSAGES
DROP INDEX IF EXISTS messages_session_paging_idx;
DROP INDEX IF EXISTS messages_dialog_paging_idx;
DROP INDEX IF EXISTS messages_sender_time_idx;

-- DIALOG MEMBERS (endi kerak emas)
DROP INDEX IF EXISTS dialog_members_user_idx;

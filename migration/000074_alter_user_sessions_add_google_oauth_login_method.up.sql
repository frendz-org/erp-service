-- Add GOOGLE_OAUTH to user_sessions.login_method CHECK constraint
ALTER TABLE user_sessions DROP CONSTRAINT IF EXISTS user_sessions_login_method_check;
ALTER TABLE user_sessions ADD CONSTRAINT user_sessions_login_method_check
    CHECK (login_method IN ('EMAIL_OTP', 'GOOGLE_OAUTH'));

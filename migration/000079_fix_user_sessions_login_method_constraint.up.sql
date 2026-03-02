-- Drop legacy constraint that only allows EMAIL_OTP
ALTER TABLE user_sessions DROP CONSTRAINT IF EXISTS chk_user_sessions_login_method;

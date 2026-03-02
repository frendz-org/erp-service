-- Re-add legacy constraint (EMAIL_OTP only)
ALTER TABLE user_sessions ADD CONSTRAINT chk_user_sessions_login_method
    CHECK (login_method IN ('EMAIL_OTP'));

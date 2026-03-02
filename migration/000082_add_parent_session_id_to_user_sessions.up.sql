ALTER TABLE user_sessions ADD COLUMN parent_session_id UUID;

ALTER TABLE user_sessions ADD CONSTRAINT fk_user_sessions_parent_session
    FOREIGN KEY (parent_session_id)
    REFERENCES user_sessions(id) ON DELETE SET NULL;

CREATE INDEX idx_user_sessions_parent_session_id
    ON user_sessions(parent_session_id)
    WHERE parent_session_id IS NOT NULL;

COMMENT ON COLUMN user_sessions.parent_session_id IS 'Links to the source session when created via transfer token. NULL for primary login sessions.';

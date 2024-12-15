CREATE TABLE user_bans (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    reason TEXT NOT NULL,
    expires_at TIMESTAMP NULL,  -- NULL 表示永久封禁
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT UNSIGNED NOT NULL,  -- 记录是哪个管理员进行的封禁
    is_active BOOLEAN DEFAULT TRUE,  -- 用于提前解除封禁
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX idx_user_bans_user_id ON user_bans(user_id);
CREATE INDEX idx_user_bans_is_active ON user_bans(is_active); 

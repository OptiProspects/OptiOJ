CREATE TABLE user_logins (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    login_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45) NOT NULL,  -- IPv6 地址最长 45 字符
    user_agent TEXT,                  -- 记录用户浏览器信息
    login_status VARCHAR(20) NOT NULL, -- success, failed, blocked 等
    fail_reason TEXT,                 -- 登录失败原因
    location VARCHAR(100),            -- 登录地理位置（可选）
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_user_logins_user_id ON user_logins(user_id);
CREATE INDEX idx_user_logins_login_time ON user_logins(login_time);
CREATE INDEX idx_user_logins_ip_address ON user_logins(ip_address); 

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_id BIGINT UNSIGNED,  -- NULL表示系统消息
    receiver_id BIGINT UNSIGNED NOT NULL,
    type VARCHAR(50) NOT NULL,  -- 消息类型：system(系统消息), team_application(团队申请), team_invitation(团队邀请)等
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    is_processed BOOLEAN DEFAULT false,  -- 是否已处理（用于申请类消息）
    application_id BIGINT UNSIGNED,      -- 关联的申请ID
    read_at TIMESTAMP NULL,     -- 消息阅读时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id)
);

CREATE TABLE team_applications (
    id SERIAL PRIMARY KEY,
    team_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    message TEXT,  -- 申请消息
    status ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE KEY (team_id, user_id, status)  -- 防止重复申请
);

-- 创建索引
CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_team_applications_team_id ON team_applications(team_id);
CREATE INDEX idx_team_applications_user_id ON team_applications(user_id);
CREATE INDEX idx_team_applications_status ON team_applications(status);

CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    avatar VARCHAR(255),
    created_by BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE team_members (
    team_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    role ENUM('owner', 'admin', 'member') NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (team_id, user_id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE team_assignments (
    id SERIAL PRIMARY KEY,
    team_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_by BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE team_assignment_problems (
    assignment_id BIGINT UNSIGNED NOT NULL,
    problem_id BIGINT UNSIGNED NOT NULL,
    order_index INT NOT NULL,
    score INT NOT NULL DEFAULT 100,
    PRIMARY KEY (assignment_id, problem_id),
    FOREIGN KEY (assignment_id) REFERENCES team_assignments(id),
    FOREIGN KEY (problem_id) REFERENCES problems(id)
);

CREATE TABLE team_problem_lists (
    id SERIAL PRIMARY KEY,
    team_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_by BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE team_problem_list_items (
    list_id BIGINT UNSIGNED NOT NULL,
    problem_id BIGINT UNSIGNED NOT NULL,
    order_index INT NOT NULL,
    note TEXT,
    PRIMARY KEY (list_id, problem_id),
    FOREIGN KEY (list_id) REFERENCES team_problem_lists(id),
    FOREIGN KEY (problem_id) REFERENCES problems(id)
);

CREATE TABLE team_invitations (
    id SERIAL PRIMARY KEY,
    team_id BIGINT UNSIGNED NOT NULL,
    code VARCHAR(20) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_by BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (team_id, code),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE team_avatars (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    team_id BIGINT UNSIGNED NOT NULL,
    filename VARCHAR(255) NOT NULL,
    upload_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY (team_id),
    FOREIGN KEY (team_id) REFERENCES teams(id)
);

-- 创建索引
CREATE INDEX idx_team_members_user_id ON team_members(user_id);
CREATE INDEX idx_team_assignments_team_id ON team_assignments(team_id);
CREATE INDEX idx_team_problem_lists_team_id ON team_problem_lists(team_id);
CREATE INDEX idx_team_invitations_code ON team_invitations(code); 
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
    problem_type ENUM('global', 'team') NOT NULL DEFAULT 'global',
    team_problem_id BIGINT UNSIGNED NULL,
    PRIMARY KEY (assignment_id, problem_id),
    FOREIGN KEY (assignment_id) REFERENCES team_assignments(id),
    FOREIGN KEY (problem_id) REFERENCES problems(id),
    FOREIGN KEY (team_problem_id) REFERENCES team_problems(id)
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

-- 团队内名称表
CREATE TABLE team_nicknames (
    team_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    nickname VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (team_id, user_id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 团队私有题目表
CREATE TABLE team_problems (
    id SERIAL PRIMARY KEY,
    team_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    input_description TEXT,
    output_description TEXT,
    sample_cases TEXT,
    hint TEXT,
    time_limit INT NOT NULL DEFAULT 1000,
    memory_limit INT NOT NULL DEFAULT 256,
    created_by BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- 团队私有题目的测试用例表
CREATE TABLE team_problem_testcases (
    id SERIAL PRIMARY KEY,
    problem_id BIGINT UNSIGNED NOT NULL,
    input_data TEXT NOT NULL,
    output_data TEXT NOT NULL,
    is_sample BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (problem_id) REFERENCES team_problems(id)
);

-- 创建索引
CREATE INDEX idx_team_members_user_id ON team_members(user_id);
CREATE INDEX idx_team_assignments_team_id ON team_assignments(team_id);
CREATE INDEX idx_team_problem_lists_team_id ON team_problem_lists(team_id);
CREATE INDEX idx_team_invitations_code ON team_invitations(code);
CREATE INDEX idx_team_problems_team_id ON team_problems(team_id);
CREATE INDEX idx_team_problem_testcases_problem_id ON team_problem_testcases(problem_id);

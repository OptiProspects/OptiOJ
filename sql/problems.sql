CREATE TABLE problems (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    input_description TEXT,
    output_description TEXT,
    samples TEXT,  -- JSON格式存储样例输入输出
    hint TEXT,
    source VARCHAR(255),
    difficulty_system ENUM('normal', 'oi') NOT NULL DEFAULT 'normal',
    difficulty VARCHAR(20) NOT NULL DEFAULT 'unrated',  -- 难度等级
    time_limit INT NOT NULL DEFAULT 1000,  -- 单位：毫秒
    memory_limit INT NOT NULL DEFAULT 256,  -- 单位：MB
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_by BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE problem_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    parent_id BIGINT UNSIGNED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES problem_categories(id)
);

CREATE TABLE problem_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL UNIQUE,
    color VARCHAR(7) DEFAULT '#000000',  -- 标签颜色，十六进制格式
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE problem_category_relations (
    problem_id BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (problem_id, category_id),
    FOREIGN KEY (problem_id) REFERENCES problems(id),
    FOREIGN KEY (category_id) REFERENCES problem_categories(id)
);

CREATE TABLE problem_tag_relations (
    problem_id BIGINT UNSIGNED NOT NULL,
    tag_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (problem_id, tag_id),
    FOREIGN KEY (problem_id) REFERENCES problems(id),
    FOREIGN KEY (tag_id) REFERENCES problem_tags(id)
);

CREATE TABLE test_cases (
    id SERIAL PRIMARY KEY,
    problem_id BIGINT UNSIGNED NOT NULL,
    input_file VARCHAR(255) NOT NULL,  -- 输入文件路径
    output_file VARCHAR(255) NOT NULL,  -- 输出文件路径
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (problem_id) REFERENCES problems(id)
);

-- 创建索引
CREATE INDEX idx_problems_difficulty ON problems(difficulty);
CREATE INDEX idx_problems_is_public ON problems(is_public);
CREATE INDEX idx_problems_created_by ON problems(created_by);
CREATE INDEX idx_test_cases_problem_id ON test_cases(problem_id);

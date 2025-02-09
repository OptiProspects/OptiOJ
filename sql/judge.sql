CREATE TABLE submissions (
    id SERIAL PRIMARY KEY,
    problem_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    language VARCHAR(20) NOT NULL,  -- 编程语言：c, cpp, java, python, etc.
    code TEXT NOT NULL,             -- 提交的代码
    status VARCHAR(20) NOT NULL,    -- 判题状态：pending, judging, accepted, wrong_answer, etc.
    time_used INT,                  -- 运行时间（毫秒）
    memory_used INT,                -- 内存使用（KB）
    error_message TEXT,             -- 错误信息
    assignment_id BIGINT UNSIGNED,  -- 作业ID，为空表示非作业提交
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (problem_id) REFERENCES problems(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (assignment_id) REFERENCES team_assignments(id)
);

CREATE TABLE judge_results (
    id SERIAL PRIMARY KEY,
    submission_id BIGINT UNSIGNED NOT NULL,
    test_case_id BIGINT UNSIGNED NOT NULL,
    status VARCHAR(20) NOT NULL,    -- 测试点状态：accepted, wrong_answer, time_limit_exceeded, etc.
    time_used INT,                  -- 运行时间（毫秒）
    memory_used INT,                -- 内存使用（KB）
    error_message TEXT,             -- 错误信息
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (submission_id) REFERENCES submissions(id),
    FOREIGN KEY (test_case_id) REFERENCES test_cases(id)
);

-- 创建索引
CREATE INDEX idx_submissions_problem_id ON submissions(problem_id);
CREATE INDEX idx_submissions_user_id ON submissions(user_id);
CREATE INDEX idx_submissions_status ON submissions(status);
CREATE INDEX idx_judge_results_submission_id ON judge_results(submission_id); 
CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    bio VARCHAR(255),
    gender VARCHAR(10),
    school VARCHAR(100),
    birthday TIMESTAMP,
    location VARCHAR(100),
    real_name VARCHAR(50),
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE KEY (user_id)
);
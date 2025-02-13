CREATE TABLE IF NOT EXISTS apps (
    app_id INT PRIMARY KEY,
    app_name TEXT NOT NULL UNIQUE,
    app_secret TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

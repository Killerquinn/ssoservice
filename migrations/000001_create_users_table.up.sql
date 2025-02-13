CREATE TABLE users (
    userid SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashedpassw TEXT NOT NULL,
    avatar TEXT,
    last_login TIMESTAMP,
    login_attempts INT DEFAULT 0,
    account_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE is_admin(
    userid INT REFERENCES users(userid) ON DELETE CASCADE,
    is_admin BOOLEAN NOT NULL,
    PRIMARY KEY (userid)
);

CREATE TABLE roles(
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(255) NOT NULL UNIQUE,
    role_description TEXT
);

CREATE TABLE user_roles(
    userid INT REFERENCES users(userid) ON DELETE CASCADE,
    role_id INT REFERENCES roles(role_id),
    PRIMARY KEY (userid, role_id)
);
CREATE TABLE IF NOT EXISTS permission(
    id SERIAL PRIMARY KEY,
    perm_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS role_permissions(
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permission(id) ON DELETE CASCADE
);
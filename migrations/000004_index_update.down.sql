DROP INDEX IF EXISTS idx_user_roles_user_id ON user_roles(userid);
DROP INDEX IF EXISTS idx_user_roles_role_id ON user_roles(role_id);
DROP INDEX IF EXISTS idx_users_user_id ON users(userid);
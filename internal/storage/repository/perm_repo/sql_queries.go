package permrepo

const (
	getRolePermits = `
    SELECT DISTINCT p.id AS permission_id
    FROM users u
    JOIN user_roles ur ON u.userid = ur.userid
    JOIN roles r ON ur.role_id = r.role_id
    JOIN role_permissions rp ON r.role_id = rp.role_id
    JOIN permissions p ON rp.permission_id = p.id
    WHERE u.userid = $1
    `

	appExists = `
    SELECT EXISTS(SELECT 1 FROM apps WHERE app_id = $1)
    `
)

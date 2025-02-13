package statrepo

const (
	getIsUsrBanned = `
	SELECT account_locked FROM users
	WHERE userid = $1
	`

	getRolesByID = `
	SELECT r.role_name u.username FROM users u
	JOIN user_roles ur ON u.user_id = ur.userid
	JOIN roles r ON ur.role_id = r.role_id
	WHERE u.userid = $1
	`

	getLastUserLogin = `
	SELECT last_login FROM users
	WHERE userid = $1
	`
)

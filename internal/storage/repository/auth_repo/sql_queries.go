package userrepository

const (
	createUserQuery = `
	INSERT INTO users(
    username,
    email,
    hashedpassw,
    avatar) VALUES (
	$1, $2, $3, COALESCE(NULLIF($4, ''), null)
	) RETURNING user_id, created_at, last_login`

	selectUserQuery = `
	SELECT user_id, email, hashedpassw, role
	FROM users
	WHERE email = $1
	`

	selectIsUserAdmin = `
	SELECT COALESCE(ia.is_admin, FALSE) AS is_admin
	FROM users u
	LEFT JOIN is_admin ia 
	ON u.user_id = ia.user_id
	WHERE u.user_id = $1
	`

	appSelectQuery = `
	SELECT app_id, name
	FROM apps
	WHERE app_id = $1
	`

	saveRefreshQuery = `
	INSERT INTO refresh_tokens(
	token,
	user_id,
	expires_at) VALUES (
	$1, $2, $3
	)
	)
	`
	//TODO: imlement user repository, add more queries, update protobuf, solve problem with migrator.go
)

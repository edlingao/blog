package queries

const AddUserQuery = `
	INSERT INTO users (id, username, password_hash, role)
	VALUES (:id, :username, :password_hash, :role);
`

const GetUserByUsernameQuery = `
	SELECT id, username, password_hash, role
	FROM users
	WHERE username = :username;
`

const UpdateUserQuery = `
	UPDATE users
	SET
		username = :username,
		email = :email,
		password_hash = :password_hash,
		role = :role
		updated_at = NOW()
	WHERE id = :id;
`

const DeleteUserQuery = `
	DELETE FROM users
	WHERE id = :id;
`

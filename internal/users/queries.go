package users

const (
	deleteQuery                = "DELETE FROM users WHERE id = ?"
	getAllQuery                = "SELECT id, firstname, lastname, email, age, height, active, created_at FROM users"
	updateQuery                = "UPDATE users SET firstname = ?, lastname = ?, email = ?, age = ?, height = ?, active = ?, created_at = ? WHERE id = ?"
	updateUserAgeLastnameQuery = "UPDATE users SET lastname = ?, age = ? WHERE id = ?"
	updateUserFirstnameQuery   = "UPDATE users SET firstname = ? WHERE id = ?"
	storeQuery                 = "INSERT INTO users (firstname, lastname, email, age, height, active, created_at) VALUES( ?, ?, ?, ?, ?, ?, ?)"
)

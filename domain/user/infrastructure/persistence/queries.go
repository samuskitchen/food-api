package persistence

const(

	// selectAllUser is a query that selects all rows in the user table
	selectAllUser = "SELECT id, names, last_names, email, created_at, updated_at FROM user WHERE deleted_at IS NULL;"

	// selectUserById is a query that selects a row from the user table based off of the given id.
	selectUserById = "SELECT id, names, last_names, email, created_at, updated_at FROM user WHERE id = $1;"

	// selectUserByEmailAndPassWord is a query that selects a row from the user table based off of the given id.
	selectUserByEmailAndPassWord = "SELECT id, names, last_names, email, password, created_at, updated_at FROM user WHERE email = $1 AND password = $2;"

	// insertUser is a query that inserts a new row in the user table using the values
	// given in order for names, last_names, username, email, password, created_at and updated_at.
	insertUser = "INSERT INTO user (names, last_names, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"

	// updateUser is a query that updates a row in the user table based off of id.
	// The values able to be updated are names, last_names, email and updated_at.
	updateUser = "UPDATE user SET names=$1, last_names=$2, email=$3, updated_at=$4 WHERE id=$5;"

	// deleteUser is a query that deletes a row in the user table given a id.
	deleteUser = "DELETE FROM user WHERE id=$1;"
)

package user

const(

	// selectAllUserTest is a query that selects all rows in the user table
	selectAllUserTest = "SELECT id, names, last_names, email FROM \"user\" WHERE deleted_at IS NULL ORDER BY created_at DESC;"

	// selectUserByIdTest is a query that selects a row from the user table based off of the given id.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	selectUserByIdTest = "SELECT id, names, last_names, email FROM \"user\" WHERE id \\= \\$1;"

	// selectUserByEmail is a query that selects a row from the user table based off of the given email.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	selectUserByEmailTest = "SELECT id, names, last_names, email, \"password\", created_at, updated_at FROM \"user\" WHERE email \\= \\$1;"

	// insertUserTest is a query that inserts a new row in the user table using the values
	// given in order for id, names, last_names, username, email, password, created_at and updated_at.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	insertUserTest = "INSERT INTO \"user\" \\(id, names, last_names, email, \"password\", created_at, updated_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7\\) RETURNING id, names, last_names, email;"

	// updateUserTest is a query that updates a row in the user table based off of id.
	// The values able to be updated are names, last_names, email and updated_at.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	updateUserTest = "UPDATE \"user\" SET names\\=\\$1, last_names\\=\\$2, email\\=\\$3, updated_at\\=\\$4 WHERE id\\=\\$5;"

	// deleteUserTest is a query that deletes a row in the user table given a id.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	deleteUserTest = "DELETE FROM \"user\" WHERE id\\=\\$1;"
)

package food

const(

	// selectAllFoodTest is a query that selects all rows in the food table
	selectAllFoodTest = "SELECT id, user_id, title, description, food_image FROM food WHERE deleted_at IS NULL ORDER BY created_at DESC;"

	// selectFoodByIdTest is a query that selects a row from the food table based off of the given id.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	selectFoodByIdTest = "SELECT id, user_id, title, description, food_image FROM food WHERE id \\= \\$1;"

	// selectFoodByUserIdTest is a query that selects a row from the food table based off of the given user userId.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	selectFoodByUserIdTest = "SELECT id, user_id, title, description, food_image FROM food WHERE user_id \\= \\$1;"

	// insertFoodTest is a query that inserts a new row in the user table using the values
	// given in order for id, user_id, title, description, food_image, created_at, updated_at.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	insertFoodTest = "INSERT INTO food \\(id, user_id, title, description, food_image, created_at, updated_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7\\) RETURNING id, user_id, title, description, food_image;"

	// updateFoodTest is a query that updates a row in the user table based off of id.
	// The values able to be updated are names, last_names, email and updated_at.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	updateFoodTest = "UPDATE food SET title\\=\\$1, description\\=\\$2, food_image\\=\\$3, updated_at\\=\\$4 WHERE id\\=\\$5;"

	// deleteFoodTest is a query that deletes a row in the food table given a id.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	deleteFoodTest = "DELETE FROM food WHERE id\\=\\$1;"
)

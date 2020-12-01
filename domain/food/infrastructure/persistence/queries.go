package persistence

const(

	// selectAllFood is a query that selects all rows in the food table
	selectAllFood = "SELECT id, user_id, title, description, food_image FROM food WHERE deleted_at IS NULL ORDER BY created_at DESC;"

	// selectFoodById is a query that selects a row from the food table based off of the given id.
	selectFoodById = "SELECT id, user_id, title, description, food_image FROM food WHERE id = $1;"

	// selectFoodByUserId is a query that selects a row from the food table based off of the given user userId.
	selectFoodByUserId = "SELECT id, user_id, title, description, food_image FROM food WHERE user_id = $1;"

	// insertFood is a query that inserts a new row in the user table using the values
	// given in order for id, user_id, title, description, food_image, created_at, updated_at.
	insertFood = "INSERT INTO food (id, user_id, title, description, food_image, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, user_id, title, description, food_image;"

	// updateFood is a query that updates a row in the user table based off of id.
	// The values able to be updated are names, last_names, email and updated_at.
	updateFood = "UPDATE food SET title=$1, description=$2, food_image=$3, updated_at=$4 WHERE id=$5;"

	// deleteFood is a query that deletes a row in the food table given a id.
	deleteFood = "DELETE FROM food WHERE id=$1;"
)

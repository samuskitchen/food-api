package seed

import (
	"database/sql"
	db "food-api/infrastructure/database"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"

	modelFood "food-api/domain/food/domain/model"
	modelUSer "food-api/domain/user/domain/model"
	_ "github.com/lib/pq"
)

// Open returns a new database connection for the test database.
func Open() *db.Data {
	return db.NewTest()
}

// Truncate removes all seed data from the test database.
func Truncate(dbc *sql.DB) error {
	stmt := "TRUNCATE TABLE food RESTART IDENTITY CASCADE; TRUNCATE TABLE \"user\" RESTART IDENTITY CASCADE;"

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate test database tables")
	}

	return nil
}

// UsersSeed handles seeding the user table in the database for integration tests.
func UsersSeed(dbc *sql.DB) ([]modelUSer.User, error) {
	now := time.Now()

	users := []modelUSer.User{
		{
			ID:        uuid.New().String(),
			Names:     "Daniel",
			LastNames: "De La Pava Suarez",
			Email:     "daniel.delapava@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uuid.New().String(),
			Names:     "Rebecca",
			LastNames: "Romero",
			Email:     "rebecca.romero@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for i := range users {
		query := `INSERT INTO "user" (id, names, last_names, email, "password", created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

		stmt, err := dbc.Prepare(query)
		if err != nil {
			return nil, errors.Wrap(err, "prepare user insertion")
		}

		err = users[i].HashPassword()
		if err != nil {
			return nil, errors.Wrap(err, "prepare hash password")
		}

		row := stmt.QueryRow(&users[i].ID, &users[i].Names, &users[i].LastNames, &users[i].Email, &users[i].PasswordHash, &users[i].CreatedAt, &users[i].UpdatedAt)

		if err = row.Scan(&users[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				return nil, errors.Wrap(err, "close psql statement")
			}

			return nil, errors.Wrap(err, "capture user id")
		}

		if err := stmt.Close(); err != nil {
			return nil, errors.Wrap(err, "close psql statement")
		}
	}

	return users, nil
}

// FoodsSeed handles seeding the food table in the database for integration tests.
func FoodsSeed(dbc *sql.DB, users []modelUSer.User) ([]modelFood.Food, error) {
	now := time.Now()

	foods := []modelFood.Food{
		{
			ID:          uuid.New().String(),
			UserID:      users[0].ID,
			Title:       "Title",
			Description: "Description",
			FoodImage:   "/profile-photos/food_api/309-3092053_gopher-link-transparent-cartoons-gopher-link.png",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.New().String(),
			UserID:      users[1].ID,
			Title:       "Title",
			Description: "Description",
			FoodImage:   "/profile-photos/food_api/309-3092053_gopher-link-transparent-cartoons-gopher-link.png",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	for i := range foods {
		query := `INSERT INTO food (id, user_id, title, description, food_image, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

		stmt, err := dbc.Prepare(query)
		if err != nil {
			return nil, errors.Wrap(err, "prepare food insertion")
		}

		row := stmt.QueryRow(&foods[i].ID, &foods[i].UserID, &foods[i].Title, &foods[i].Description, &foods[i].FoodImage, &foods[i].CreatedAt, &foods[i].UpdatedAt)

		if err = row.Scan(&foods[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				return nil, errors.Wrap(err, "close psql statement")
			}

			return nil, errors.Wrap(err, "capture food id")
		}

		if err := stmt.Close(); err != nil {
			return nil, errors.Wrap(err, "close psql statement")
		}
	}

	return foods, nil
}
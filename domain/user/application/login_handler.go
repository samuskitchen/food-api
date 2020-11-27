package application

import (
	repoDomain "food-api/domain/user/domain/respository"
	"food-api/domain/user/infrastructure/persistence"
	"food-api/infrastructure/database"
)

// UserRouter
type UserRouter struct {
	Repo repoDomain.UserRepository
}

// NewLoginHandler
func NewLoginHandler(db *database.Data) *UserRouter  {
	return &UserRouter{
		Repo: persistence.NewUserRepository(db),
	}
}

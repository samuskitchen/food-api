package v1

import (
	repoDomain "food-api/domain/user/domain/respository"
	"food-api/domain/user/infrastructure/persistence"
	"food-api/infrastructure/database"
)

// UserRouter
type UserRouter struct {
	Repo repoDomain.UserRepository
}

// NewUserHandler
func NewUserHandler(db *database.Data) *UserRouter  {
	return &UserRouter{
		Repo: persistence.NewUserRepository(db),
	}
}

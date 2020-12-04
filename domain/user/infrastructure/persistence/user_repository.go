package persistence

import (
	"context"
	"errors"
	"food-api/domain/user/application/v1/response"
	"food-api/domain/user/domain/model"
	repoDomain "food-api/domain/user/domain/repository"
	"food-api/infrastructure/database"
	"github.com/google/uuid"
	"strings"
)

type sqlUserRepo struct {
	Conn *database.Data
}

func NewUserRepository(Conn *database.Data) repoDomain.UserRepository {
	return &sqlUserRepo{
		Conn: Conn,
	}
}

func (sr *sqlUserRepo) GetAllUser(ctx context.Context) ([]response.UserResponse, error) {
	rows, err := sr.Conn.DB.QueryContext(ctx, selectAllUser)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []response.UserResponse
	for rows.Next() {
		var userRow response.UserResponse
		_ = rows.Scan(&userRow.ID, &userRow.Names, &userRow.LastNames, &userRow.Email)
		users = append(users, userRow)
	}

	return users, nil
}

func (sr *sqlUserRepo) GetById(ctx context.Context, id string) (response.UserResponse, error) {
	row := sr.Conn.DB.QueryRowContext(ctx, selectUserById, id)

	var userScan response.UserResponse
	err := row.Scan(&userScan.ID, &userScan.Names, &userScan.LastNames, &userScan.Email)
	if err != nil {
		return response.UserResponse{}, err
	}

	return userScan, nil
}

func (sr *sqlUserRepo) CreateUser(ctx context.Context, user *model.User) (*response.UserResponse, error) {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, insertUser)
	if err != nil {
		return &response.UserResponse{}, err
	}

	defer stmt.Close()
	if strings.TrimSpace(user.ID) == "" {
		user.ID = uuid.New().String()
	}

	row := stmt.QueryRowContext(ctx, &user.ID, &user.Names, &user.LastNames, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	userResult := response.UserResponse{}
	err = row.Scan(&userResult.ID, &userResult.Names, &userResult.LastNames, &userResult.Email)
	if err != nil {
		return &response.UserResponse{}, err
	}

	return &userResult, nil
}

func (sr *sqlUserRepo) UpdateUser(ctx context.Context, id string, user model.User) error {
	stmt, err := sr.Conn.DB.PrepareContext(ctx, updateUser)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.Names, user.LastNames, user.Email, user.UpdatedAt, id)
	if err != nil {
		return err
	}

	return err
}

func (sr *sqlUserRepo) GetUserByEmailAndPassword(ctx context.Context, user *model.User) (*response.UserResponse, error) {
	row := sr.Conn.DB.QueryRowContext(ctx, selectUserByEmail, user.Email)

	userScan := model.User{}
	err := row.Scan(&userScan.ID, &userScan.Names, &userScan.LastNames, &userScan.Email, &userScan.PasswordHash, &userScan.CreatedAt, &userScan.UpdatedAt)
	if err != nil {
		return &response.UserResponse{}, err
	}

	validate := userScan.PasswordMatch(user.Password)
	if !validate {
		return &response.UserResponse{}, errors.New("password does not match")
	}

	userResponse := response.UserResponse{
		ID:        userScan.ID,
		Names:     userScan.Names,
		LastNames: userScan.LastNames,
		Email:     userScan.Email,
	}

	return &userResponse, nil
}

package persistence

import (
	"context"
	"food-api/domain/user/application/v1/response"
	"food-api/domain/user/domain/model"
	repoDomain "food-api/domain/user/domain/respository"
	"food-api/infrastructure/database"
	"time"
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

func (sr *sqlUserRepo) CreateUser(ctx context.Context, user model.User) (response.UserResponse, error) {
	now := time.Now() //.Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	stmt, err := sr.Conn.DB.PrepareContext(ctx, insertUser)
	if err != nil {
		return response.UserResponse{}, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, user.Names, user.LastNames, user.Email, user.Password, now, now)

	userResponse := response.UserResponse{}
	err = row.Scan(&userResponse.ID, &userResponse.Names, &userResponse.LastNames, &userResponse.Email)
	if err != nil {
		return response.UserResponse{}, err
	}

	return userResponse, nil
}

func (sr *sqlUserRepo) UpdateUser(ctx context.Context, id string, user model.User) error {
	now := time.Now() //.Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	stmt, err := sr.Conn.DB.PrepareContext(ctx, updateUser)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.Names, user.LastNames, user.Email, now, id)
	if err != nil {
		return err
	}

	return err
}

func (sr *sqlUserRepo) GetUserByEmailAndPassword(ctx context.Context, user *model.User) (*response.UserResponse, error) {
	row := sr.Conn.DB.QueryRowContext(ctx, selectUserByEmailAndPassWord, user.Email, user.Password)

	var userScan response.UserResponse
	err := row.Scan(&userScan.ID, &userScan.Names, &userScan.LastNames, &userScan.Email)
	if err != nil {
		return &response.UserResponse{}, err
	}

	return &userScan, nil
}
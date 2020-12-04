package user

import (
	"context"
	"database/sql"
	responseUser "food-api/domain/user/application/v1/response"
	"food-api/domain/user/domain/model"
	"food-api/domain/user/domain/repository"
	"food-api/domain/user/infrastructure/persistence"
	"food-api/infrastructure/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

// represent the repository
var (
	dbMockUsers        *sql.DB
	connMockUser       database.Data
	userRepositoryMock repository.UserRepository
)

// NewMockUser initialize mock connection to database
func NewMockUser() sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbMockUsers = db
	connMockUser = database.Data{
		DB: dbMockUsers,
	}

	userRepositoryMock = persistence.NewUserRepository(&connMockUser)

	return mock
}

// CloseMockUser Close attaches the provider and close the connection
func CloseMockUser() {
	err := dbMockUsers.Close()
	if err != nil {
		log.Println("Error close database test")
	}
}

// dataUser is data for test
func dataUser() []model.User {
	now := time.Now()//.Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	return []model.User{
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
}

// dataUserResponse is data for test
func dataUserResponse() []responseUser.UserResponse {

	return []responseUser.UserResponse{
		{
			ID:        uuid.New().String(),
			Names:     "Daniel",
			LastNames: "De La Pava Suarez",
			Email:     "daniel.delapava@jikkosoft.com",
		},
		{
			ID:        uuid.New().String(),
			Names:     "Rebecca",
			LastNames: "Romero",
			Email:     "rebecca.romero@jikkosoft.com",
		},
	}
}

func Test_sqlUserRepo_GetAllUser(t *testing.T) {

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		mock.ExpectQuery("SELECT 1 FROM user")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		users, err := userRepositoryMock.GetAllUser(ctx)
		assert.Error(tt, err)
		assert.Nil(tt, users)
	})

	t.Run("Get All User Successful", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		usersData := dataUserResponse()
		rows := sqlmock.NewRows([]string{"id", "names", "last_names", "email"}).
			AddRow(usersData[0].ID, usersData[0].Names, usersData[0].LastNames, usersData[0].Email).
			AddRow(usersData[1].ID, usersData[1].Names, usersData[1].LastNames, usersData[1].Email)

		mock.ExpectQuery(selectAllUserTest).WillReturnRows(rows)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		users, err := userRepositoryMock.GetAllUser(ctx)
		assert.NotEmpty(tt, users)
		assert.NoError(tt, err)
		assert.Len(tt, users, 2)
	})

}

func Test_sqlUserRepo_GetById(t *testing.T) {
	userTest := dataUserResponse()[0]

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		row := sqlmock.NewRows([]string{"id", "names", "last_names", "email"}).
			AddRow(userTest.ID, userTest.Names, userTest.LastNames, userTest.Email)

		mock.ExpectQuery(selectUserByIdTest).WithArgs(nil).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.GetById(ctx, "")
		assert.Error(tt, err)
		assert.NotNil(tt, userResult)
	})

	t.Run("Get User By Id Successful", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		row := sqlmock.NewRows([]string{"id", "names", "last_names", "email"}).
			AddRow(userTest.ID, userTest.Names, userTest.LastNames, userTest.Email)

		mock.ExpectQuery(selectUserByIdTest).WithArgs(userTest.ID).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.GetById(ctx, userTest.ID)
		assert.NoError(tt, err)
		assert.NotNil(tt, userResult)
	})
}

func Test_sqlUserRepo_GetUserByEmailAndPassword(t *testing.T) {
	userTest := dataUser()[0]

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		row := sqlmock.NewRows([]string{"id", "names", "last_names", "email"}).
			AddRow(userTest.ID, userTest.Names, userTest.LastNames, userTest.Email)

		mock.ExpectQuery(selectUserByEmailTest).WithArgs(nil).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.GetUserByEmailAndPassword(ctx, &dataUser()[0])
		assert.Error(tt, err)
		assert.NotNil(tt, userResult)
	})

	t.Run("Error Password Not Match", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		row := sqlmock.NewRows([]string{"id", "names", "last_names", "email", "password", "created_at", "updated_at"}).
			AddRow(userTest.ID, userTest.Names, userTest.LastNames, userTest.Email, userTest.PasswordHash, userTest.CreatedAt, userTest.UpdatedAt)

		mock.ExpectQuery(selectUserByEmailTest).WithArgs(userTest.Email).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.GetUserByEmailAndPassword(ctx, &dataUser()[0])
		assert.Error(tt, err)
		assert.NotNil(tt, userResult)
	})

	t.Run("Get User By Email Successful", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		err := userTest.HashPassword()
		if err != nil {
			assert.Error(tt, err)
		}

		row := sqlmock.NewRows([]string{"id", "names", "last_names", "email", "password", "created_at", "updated_at"}).
			AddRow(userTest.ID, userTest.Names, userTest.LastNames, userTest.Email, userTest.PasswordHash, userTest.CreatedAt, userTest.UpdatedAt)

		mock.ExpectQuery(selectUserByEmailTest).WithArgs(userTest.Email).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.GetUserByEmailAndPassword(ctx, &dataUser()[0])
		assert.NoError(tt, err)
		assert.NotNil(tt, userResult)
	})

}

func Test_sqlUserRepo_CreateUser(t *testing.T) {

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare("insertUserTest")
		prep.ExpectExec().
			WithArgs(dataUser()[0].Names, dataUser()[0].LastNames, dataUser()[0].Email, dataUser()[0].Password, dataUser()[0].CreatedAt, dataUser()[0].UpdatedAt).
			WillReturnResult(sqlmock.NewResult(0, 0))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.CreateUser(ctx, &dataUser()[0])
		assert.Error(tt, err)
		assert.NotNil(tt, userResult)
	})

	t.Run("Error Scan Row", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare(insertUserTest)
		prep.ExpectQuery().
			WithArgs(dataUser()[0].Names, dataUser()[0].LastNames, dataUser()[0].Email, dataUser()[0].PasswordHash, dataUser()[0].CreatedAt, dataUser()[0].UpdatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"first_name"}).AddRow("Error"))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.CreateUser(ctx, &dataUser()[0])
		assert.Error(tt, err)
		assert.NotNil(tt, userResult)
	})

	t.Run("Create User Successful", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		dataTest := dataUser()[0]
		err := dataTest.HashPassword()
		if err != nil {
			assert.Error(tt, err)
		}

		prep := mock.ExpectPrepare(insertUserTest)
		prep.ExpectQuery().
			WithArgs(dataTest.ID, dataTest.Names, dataTest.LastNames, dataTest.Email, dataTest.PasswordHash, dataTest.CreatedAt, dataTest.UpdatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"id", "names", "last_names", "email"}).
				AddRow(dataTest.ID, dataTest.Names, dataTest.LastNames, dataTest.Email))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userResult, err := userRepositoryMock.CreateUser(ctx, &dataTest)
		assert.NoError(tt, err)
		assert.NotNil(tt, userResult)
	})

}

func Test_sqlUserRepo_UpdateUser(t *testing.T) {

	t.Run("Error SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare("updateUserTest")
		prep.ExpectExec().
			WithArgs(dataUser()[0].Names, dataUser()[0].LastNames, dataUser()[0].Email, dataUser()[0].UpdatedAt, dataUser()[0].ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.UpdateUser(ctx, dataUser()[0].ID, dataUser()[0])
		assert.Error(tt, err)
	})

	t.Run("Error Statement SQL", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		prep := mock.ExpectPrepare(updateUserTest)
		prep.ExpectExec().
			WithArgs(dataUser()[0].Names, dataUser()[0].LastNames, dataUser()[0].Email, dataUser()[0].UpdatedAt, nil).
			WillReturnResult(sqlmock.NewResult(1, 1))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.UpdateUser(ctx, dataUser()[0].ID, dataUser()[0])
		assert.Error(tt, err)
	})

	t.Run("Update User Successful", func(tt *testing.T) {
		mock := NewMockUser()
		defer func() {
			CloseMockUser()
		}()

		dataTest := dataUser()[0]
		prep := mock.ExpectPrepare(updateUserTest)
		prep.ExpectExec().
			WithArgs(dataTest.Names, dataTest.LastNames, dataTest.Email, dataTest.UpdatedAt, dataTest.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := userRepositoryMock.UpdateUser(ctx, dataTest.ID, dataTest)
		assert.NoError(tt, err)
	})

}

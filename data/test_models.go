package data

import (
	"gorm.io/gorm"
)

type PostgresTestRepository struct {
	Conn *gorm.DB
}

func NewPostgresTestRepository(pool *gorm.DB) *PostgresTestRepository {
	return &PostgresTestRepository{
		Conn: pool,
	}
}

func (tr *PostgresTestRepository) GetByUsername(username string) (*User, error) {
	user := User{
		Username:       "test-username",
		Password:       "test-password",
		Role:           "member",
		OrganizationID: "test-org-1",
	}
	user.ID = "random-test-id"
	return &user, nil
}

func (tr *PostgresTestRepository) GetByID(id string) (*User, error) {
	user := User{
		Username:       "test-username",
		Password:       "test-password",
		Role:           "admin",
		OrganizationID: "test-org-1",
	}
	user.ID = id
	return &user, nil
}

func (tr *PostgresTestRepository) Insert(user User) error {
	return nil
}

func (tr *PostgresTestRepository) Delete(user User) error {
	return nil
}

func (tr *PostgresTestRepository) PasswordMatch(plainTextPassword string, user User) (bool, error) {
	return true, nil
}

func (tr *PostgresTestRepository) GetAllOtherUsersInOrg(user User) ([]User, error) {
	users := []User{}
	return users, nil
}

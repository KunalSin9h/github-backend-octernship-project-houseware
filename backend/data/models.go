package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

const dbQueryTimeout = time.Second * 2

type Models struct {
	User         User
	Organization Organization
}

type GormModel struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Organization struct {
	GormModel
	Name  string
	Desc  string
	Users []User `gorm:"many2many:organization_users;"`
}

type User struct {
	GormModel
	FirstName      string
	LastName       string
	Email          string
	Username       string
	Password       string
	Role           string
	OrganizationID string
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
}

func New(dbPood *gorm.DB) Models {
	db = dbPood

	db.AutoMigrate(&User{})

	return Models{
		User:         User{},
		Organization: Organization{},
	}
}

func (u *User) PasswordMatch(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainTextPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// Invalid Password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

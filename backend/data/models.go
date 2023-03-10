package data

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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
	ID        string         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type Organization struct {
	GormModel
	Name  string `json:"name" gorm:"unique"`
	Desc  string `json:"desc"`
	Users []User `json:"users" gorm:"many2many:organization_users;"`
}

type User struct {
	GormModel
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	Email          string       `json:"email" gorm:"unique"`
	Username       string       `json:"username" gorm:"unique"`
	Password       string       `json:"-"`
	Role           string       `json:"role"`
	OrganizationID string       `json:"organization_id"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
}

func New(dbPool *gorm.DB) Models {
	db = dbPool

	db.AutoMigrate(&User{})

	return Models{
		User:         User{},
		Organization: Organization{},
	}
}

func (u *User) GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	var user User
	err := db.WithContext(ctx).Model(User{}).Find(&user, "username = ?", username).Error

	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func (u *User) Insert(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	if err != nil {
		return err
	}

	user.ID = uuid.NewString()
	user.Password = string(hashPassword)

	if err != nil {
		return err
	}

	err = db.WithContext(ctx).Create(&user).Error

	if err != nil {
		return err
	}

	return nil
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

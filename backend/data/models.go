package data

import (
	"context"
	"errors"
	"log"
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
	ID        string    `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Organization struct {
	GormModel
	Name  string `json:"name" gorm:"unique"`
	Users []User
}

type User struct {
	GormModel
	Username       string `json:"username" gorm:"not null;unique"`
	Password       string `json:"-"`
	Role           string `json:"role" gorm:"not null"`
	OrganizationID string `json:"organization_id" gorm:"not null"`
}

func New(dbPool *gorm.DB) Models {
	db = dbPool

	db.AutoMigrate(&Organization{}, &User{})
	populateDatabase()

	return Models{
		User:         User{},
		Organization: Organization{},
	}
}

func (org *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	org.ID = uuid.NewString()
	return nil
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.NewString()
	return nil
}

func (u *User) GetByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	var user User
	err := db.WithContext(ctx).Model(&User{}).Find(&user, "username = ?", username).Error

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

	err = db.WithContext(ctx).Create(&user).Error

	if err != nil {
		return err
	}

	err = db.WithContext(ctx).Model(&user).Association("Organizations").Append(&Organization{})

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

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		log.Fatal("@MODELS Failed to hash password")
	}

	return string(hash)
}

func populateDatabase() {

	orgs := []Organization{
		{Name: "Apple"},
		{Name: "Google"},
	}

	for _, org := range orgs {
		err := db.Model(&Organization{}).Create(&org).Error
		if err != nil {
			log.Fatal("@MODELS Failed to populate Organization table")
		}
	}

	var appleOrg, googleOrg Organization
	db.Model(&Organization{}).First(&appleOrg, "Name = ?", "Apple")
	db.Model(&Organization{}).First(&googleOrg, "Name = ?", "Google")

	users := []User{
		{Username: "user1", Password: hashPassword("user1"), Role: "admin", OrganizationID: appleOrg.ID},
		{Username: "user2", Password: hashPassword("user2"), Role: "member", OrganizationID: appleOrg.ID},
		{Username: "user3", Password: hashPassword("user3"), Role: "admin", OrganizationID: googleOrg.ID},
		{Username: "user4", Password: hashPassword("user4"), Role: "member", OrganizationID: googleOrg.ID},
	}

	for _, user := range users {
		err := db.Model(&User{}).Create(&user).Error
		if err != nil {
			log.Fatal("@MODELS Failed to populate User table")
		}
	}

	log.Println("@MODELS Successfully populated database")
}

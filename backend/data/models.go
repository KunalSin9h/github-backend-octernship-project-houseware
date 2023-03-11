package data

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
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

func (user *User) BeforeDelete(tx *gorm.DB) (err error) {
	if user.Role == "admin" {
		return errors.New("admin user not allowed to delete")
	}
	return
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

func (u *User) GetByID(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	var user User
	err := db.WithContext(ctx).Model(&User{}).Find(&user, "id = ?", id).Error

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

	user.Password = string(hashPassword)

	err = db.WithContext(ctx).Create(&user).Error

	if err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	err := db.WithContext(ctx).Model(&User{}).Delete(&u).Error

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

func (u *User) GetAllOtherUsersInOrg() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	var users []User

	err := db.WithContext(ctx).Model(&User{}).Find(&users, "organization_id = ? and id != ?", u.OrganizationID, u.ID).Error

	if err != nil {
		return []User{}, err
	}

	return users, nil
}

func ConnectDatabase(DSN string) *gorm.DB {
	numberOfTry := 0
	numberOfTryLimit := 5

	for {
		db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

		if err != nil {
			numberOfTry++
			log.Printf("@MAIN Trying to connect to Postgres Database...[%d/%d]", numberOfTry, numberOfTryLimit)
		} else {
			log.Println("@MAIN Successfully Connected to Postgres Database")
			return db
		}

		if numberOfTry >= numberOfTryLimit {
			log.Fatal("@MAIN Failed to Connect to Postgres Database")
		}

		holdTime := numberOfTry * numberOfTry // numberOfTry ^ 2
		log.Printf("@MAIN Retrying to connect in %d sec", holdTime)
		time.Sleep(time.Duration(holdTime) * time.Second)
	}
}

func populateDatabase() {

	db.Exec("TRUNCATE users, organizations")

	orgs := []Organization{
		{Name: "ORG-1"},
		{Name: "ORG-2"},
	}

	for _, org := range orgs {
		err := db.Model(&Organization{}).Create(&org).Error
		if err != nil {
			log.Fatal("@MODELS Failed to populate Organization table")
		}
	}

	var org_1, org_2 Organization
	db.Model(&Organization{}).First(&org_1, "Name = ?", "ORG-1")
	db.Model(&Organization{}).First(&org_2, "Name = ?", "ORG-2")

	password, err := bcrypt.GenerateFromPassword([]byte("password"), 12)

	if err != nil {
		log.Fatal("@MODELS Failed to hash password for dummy data")
	}

	users := []User{
		// ORG-1 Members
		{Username: "User-1", Password: string(password), Role: "admin", OrganizationID: org_1.ID},
		{Username: "User-2", Password: string(password), Role: "admin", OrganizationID: org_1.ID},
		{Username: "User-3", Password: string(password), Role: "member", OrganizationID: org_1.ID},
		{Username: "User-4", Password: string(password), Role: "member", OrganizationID: org_1.ID},
		// ORG-2 Members
		{Username: "User-5", Password: string(password), Role: "admin", OrganizationID: org_2.ID},
		{Username: "User-6", Password: string(password), Role: "admin", OrganizationID: org_2.ID},
		{Username: "User-7", Password: string(password), Role: "member", OrganizationID: org_2.ID},
		{Username: "User-8", Password: string(password), Role: "member", OrganizationID: org_2.ID},
	}

	for _, user := range users {
		err := db.Model(&User{}).Create(&user).Error
		if err != nil {
			log.Fatal("@MODELS Failed to populate User table")
		}
	}

	log.Println("@MODELS Successfully populated database")
}

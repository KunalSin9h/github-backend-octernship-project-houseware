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

type PostgresRepository struct {
	Conn *gorm.DB
}

func NewPostgresRepository(pool *gorm.DB) *PostgresRepository {
	db = pool
	return &PostgresRepository{
		Conn: pool,
	}
}

// type Models struct {
// 	User         User
// 	Organization Organization
// }

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

/*
data.New is a function that takes a pointer to a gorm.DB object and returns a Models struct.
It sets the package level variable db to the gorm.DB object passed to it.
*/
// func New(dbPool *gorm.DB) Models {
// 	db = dbPool

// 	db.AutoMigrate(&Organization{}, &User{})
// 	populateDatabase()

// 	return Models{
// 		User:         User{},
// 		Organization: Organization{},
// 	}
// }

/*
=====================
HOOKS
=====================
*/
// BeforeCreate hook is used to generate a UUID for the ID field of the Organization struct
func (org *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	org.ID = uuid.NewString()
	return nil
}

// BeforeCreate hook is used to generate a UUID for the ID field of the User struct
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.NewString()
	return nil
}

// BeforeDelete hook is used to prevent the deletion of the admin user
func (user *User) BeforeDelete(tx *gorm.DB) (err error) {
	if user.Role == "admin" {
		return errors.New("admin user not allowed to delete")
	}
	return
}

// =====================================================

/*
GetByUsername is a method that takes a username and returns a User struct and an error.
*/
func (u *PostgresRepository) GetByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	var user User
	err := db.WithContext(ctx).Model(&User{}).Find(&user, "username = ?", username).Error

	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

/*
GetById is a method that takes an id and returns a User struct and an error.
*/
func (u *PostgresRepository) GetByID(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	var user User
	err := db.WithContext(ctx).Model(&User{}).Find(&user, "id = ?", id).Error

	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

/*
Insert is a method that inserts a User struct into the database and returns an error.
*/
func (u *PostgresRepository) Insert(user User) error {
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

/*
Delete is a method that deletes a User struct from the database and returns an error.
*/
func (u *PostgresRepository) Delete(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	err := db.WithContext(ctx).Model(&User{}).Delete(&user).Error

	if err != nil {
		return err
	}

	return nil
}

/*
PasswordMatch is a method that takes a plain text password and matches it with hash password and returns a boolean and an error.
*/
func (u *PostgresRepository) PasswordMatch(plainTextPassword string, user User) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainTextPassword))

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

/*
GetAllUsersInOrg is a method that returns all other users from the same organization
*/
func (u *PostgresRepository) GetAllOtherUsersInOrg(user User) ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbQueryTimeout)
	defer cancel()

	var users []User

	err := db.WithContext(ctx).Model(&User{}).Find(&users, "organization_id = ? and id != ?", user.OrganizationID, user.ID).Error

	if err != nil {
		return []User{}, err
	}

	return users, nil
}

/*
ConnectDatabase is used to connect to database
*/
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

/*
PopulateDatabase is used to populate the database with dummy data
You can see the visual representation of the data in the README.md file
*/
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

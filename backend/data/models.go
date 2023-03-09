package data

import (
	"time"

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

type User struct {
	GormModel
	FirstName string
	LastName  string
	Email     string
	Username  string
	Password  string
	Role      string
	Org       Organization
}

type Organization struct {
	GormModel
	Name  string
	Users []User
}

func New(dbPood *gorm.DB) Models {
	db = dbPood

	return Models{
		User:         User{},
		Organization: Organization{},
	}
}

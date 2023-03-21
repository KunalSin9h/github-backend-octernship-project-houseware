package data

/*
	Repository Method to make our database testable

	PostgresRepository----
	                      \__ Repository
						  /
	TestRepository-------
*/

type Repository interface {
	GetByUsername(username string) (*User, error)
	GetByID(id string) (*User, error)
	GetAllOtherUsersInOrg(user User) ([]User, error)
	Insert(user User) error
	Delete(user User) error
	PasswordMatch(plainTextPassword string, user User) (bool, error)
}

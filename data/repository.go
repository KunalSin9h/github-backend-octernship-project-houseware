package data

/*
	Repository Method to make our handlers testable by mocking database.

	PostgresRepository ----------
	                             \__ Repository
						         /
	PostgresTestRepository-------

	PostgresTestRepository will also define these methods but they will be dummy and return
	data without query the database
*/

type Repository interface {
	GetByUsername(username string) (*User, error)
	GetByID(id string) (*User, error)
	GetAllOtherUsersInOrg(user User) ([]User, error)
	Insert(user User) error
	Delete(user User) error
	PasswordMatch(plainTextPassword string, user User) (bool, error)
}

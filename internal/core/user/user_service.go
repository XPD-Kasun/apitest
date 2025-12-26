package user

type UserService interface {
	// Logs a user in with username, password combination.
	// Returns a JWT token.
	Login(username, password string) (string, error)

	Logout()

	CreateNew(user *AppUser) error

	GetById(id int) (*AppUser, error)

	GetByIds(id ...int) ([]*AppUser, error)

	GetByPage(cursor int, limit int) ([]*AppUser, error)
}

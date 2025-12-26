package user

import (
	"apitest/internal/core/common/baserepo"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userRepo AppUserRepo
	jwtKey   string
}

func NewUserServiceImpl(repo AppUserRepo, jwtKey string) *UserServiceImpl {
	return &UserServiceImpl{userRepo: repo, jwtKey: jwtKey}
}

// GetById implements [UserService].
func (us *UserServiceImpl) GetById(id int) (*AppUser, error) {
	return us.userRepo.GetById(id)
}

// GetByIds implements [UserService].
func (us *UserServiceImpl) GetByIds(ids ...int) ([]*AppUser, error) {
	return us.userRepo.GetByIds(ids...)
}

func (u *UserServiceImpl) GetByPage(cursor int, limit int) ([]*AppUser, error) {
	result, err := u.userRepo.GetByPage(&baserepo.PaginatedFilter[int]{
		Cursor: cursor,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (us *UserServiceImpl) Login(username, password string) (string, error) {
	user, err := us.userRepo.GetByUserName(username)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	pass, err := us.userRepo.GetPasswordHash(user.Id)
	if err != nil {
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(pass), []byte(password)) != nil {
		log.Info().Msgf("Login attempt for user %s with id=%d with password '%s' ", user.UserName, user.Id, password)
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "XPD.org",
		"sub": strconv.Itoa(user.Id),
	})

	key := os.Getenv("jwtKey")

	if key == "" {
		return "", ErrNoJWTKeyFound
	}

	st, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return st, nil
}

func (us *UserServiceImpl) Logout() {

}

func (us *UserServiceImpl) CreateNew(user *AppUser) error {

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(pass)
	if err := us.userRepo.Insert(user); err != nil {
		return err
	}

	return nil
}

// vs implementer
//var s UserService = &UserServiceImpl{}

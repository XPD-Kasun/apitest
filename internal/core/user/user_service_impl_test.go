package user

import (
	"apitest/internal/core/common/baserepo"
	"errors"
	"strconv"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

var user1 = AppUser{
	Id:        1,
	UserName:  "john",
	Email:     "john@yahoo.com",
	Password:  "johnfn",
	Firstname: "John",
	Lastname:  "McLaren",
}

var user2 = AppUser{
	Id:        2,
	UserName:  "kamal",
	Email:     "kamal@gmail.com",
	Password:  "kamalfn2#",
	Firstname: "Kamal",
	Lastname:  "Godagama",
}

var jwtKey = "jdsfoi34"

type appUserRepoStub struct {
	users         map[string]*AppUser
	errNotExist   string
	errUserExists string
}

// var s AppUserRepo = appUserRepoStub{}

// GetByPage implements [AppUserRepo].
func (a appUserRepoStub) GetByPage(filter *baserepo.PaginatedFilter[int]) (*baserepo.PaginatedResult[*AppUser, int], error) {
	cursor := filter.Cursor
	limit := filter.Limit
	ids := make([]int, 0, limit)

	for i := 0; i < limit; i++ {
		ids = append(ids, i+cursor)
	}

	items, err := a.GetByIds(ids...)
	if err != nil {
		return nil, err
	}

	_, err = a.GetById(cursor + limit)

	result := &baserepo.PaginatedResult[*AppUser, int]{
		Items:      items,
		HasMore:    err == nil,
		NextCursor: cursor + limit,
	}

	return result, nil
}

// GetById implements [AppUserRepo].
func (a appUserRepoStub) GetById(id int) (*AppUser, error) {
	for _, u := range a.users {
		if u.Id == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

// GetByIds implements [AppUserRepo].
func (a appUserRepoStub) GetByIds(ids ...int) ([]*AppUser, error) {
	var users []*AppUser
	var idMap map[int]bool = make(map[int]bool)

	for _, id := range ids {
		idMap[id] = true
	}

	for _, u := range a.users {
		if idMap[u.Id] {
			users = append(users, u)
		}
	}

	return users, nil

}

// GetByUserName implements AppUserRepo.
func (a appUserRepoStub) GetByUserName(username string) (*AppUser, error) {
	u, ok := a.users[username]
	if ok {
		return u, nil
	} else {
		return nil, errors.New(a.errNotExist)
	}
}

// GetPasswordHash implements AppUserRepo.
func (a appUserRepoStub) GetPasswordHash(id int) (string, error) {
	for _, v := range a.users {
		if v.Id == id {
			return v.Password, nil
		}
	}
	return "", errors.New(a.errNotExist)
}

// Insert implements AppUserRepo.
func (a appUserRepoStub) Insert(val *AppUser) error {
	_, ok := a.users[val.UserName]
	if ok {
		return errors.New(a.errUserExists)
	}
	a.users[val.UserName] = val
	return nil
}

func TestUserServiceImpl_CreateNewUser(t *testing.T) {

	repo := appUserRepoStub{
		errNotExist:   "user doesnt exist",
		errUserExists: "user already exist",
		users:         make(map[string]*AppUser),
	}
	userService := &UserServiceImpl{userRepo: repo}
	repo.users[user1.UserName] = &user1
	repo.users[user2.UserName] = &user2

	t.Run("InsertNewUser", func(t *testing.T) {

		err := userService.CreateNew(&AppUser{
			UserName: user1.UserName,
		})
		if err == nil {
			t.Error("inserts duplicate users")
			t.SkipNow()
		}

		newUser := AppUser{
			Id:        10,
			UserName:  "a",
			Email:     "a@a.com",
			Firstname: "A",
			Lastname:  "B",
		}

		err = userService.CreateNew(&newUser)

		if err != nil {
			t.Fatalf("adding new user fails")
		}

		sameUser, err := userService.userRepo.GetByUserName("a")

		if err != nil {
			t.Log("calling GetByUserName for newly added user returns error")
			t.Log(err)
		}

		newUser.Password, err = userService.userRepo.GetPasswordHash(newUser.Id)

		if err != nil {
			t.Log("calling GetPasswordHash for newly added user returns error")
			t.Log(err)
		}

		if *sameUser != newUser {
			t.Error("newly added user is not returned by GetByUserName")
			t.Logf("expect %v, got %v %v", newUser, sameUser, *sameUser != newUser)
		}

	})

}

func TestUserServiceImpl_Login(t *testing.T) {

	repo := appUserRepoStub{
		errNotExist:   "user doesnt exist",
		errUserExists: "user already exist",
		users:         make(map[string]*AppUser),
	}

	t.Setenv("jwtKey", jwtKey)

	userService := UserServiceImpl{userRepo: repo}
	userService.CreateNew(&user1)

	t.Run("logins with correct credentials", func(t *testing.T) {
		token, err := userService.Login(user1.UserName, user1.Password)
		if err != nil {
			t.Error("login failed with error", err)
		}
		if token == "" {
			t.Error("logged in correctly but returned empty token")
		}

		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
			return []byte(jwtKey), nil
		})
		if err != nil {
			t.Error("generated token error : ", err)
		}
		if sub, err := claims.GetSubject(); err == nil {
			if i, err2 := strconv.Atoi(sub); err2 == nil {
				if i != user1.Id {
					t.Errorf("generated token user id is wrong, expected %v, got %v", user1.Id, i)
				}
			} else {
				t.Error("sub is not a user id")
			}
		} else {
			t.Error("GetSubject failed.", err)
		}

		t.Log("token is ", token)
	})

	t.Run("cannot login with invalid password", func(t *testing.T) {

		token, err := userService.Login(user1.UserName, "gibiris")
		if err == nil {
			t.Error("logins ", err)
		}
		if token != "" {
			t.Errorf("invalid login cannot provide token %v", token)
		}
	})

	t.Run("cannot login with invalid username", func(t *testing.T) {

		token, err := userService.Login("jfosdjf", "gibiris")
		if err == nil {
			t.Error("logins ", err)
		} else {
			t.Log("invalid username login error : ", err)
		}
		if token != "" {
			t.Errorf("invalid login cannot provide token %v", token)
		}
	})
}

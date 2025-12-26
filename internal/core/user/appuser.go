package user

type AppUser struct {
	Id        int    `json:"id"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

package dto

import "apitest/internal/core/user"

type UserDTO struct {
	Id        int    `json:"id"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

func (dto *UserDTO) ToAppUser() *user.AppUser {
	var appUser = user.AppUser{
		Id:        dto.Id,
		UserName:  dto.UserName,
		Password:  dto.Password,
		Firstname: dto.Firstname,
		Lastname:  dto.Lastname,
		Email:     dto.Email,
	}
	return &appUser
}

func (dto *UserDTO) FromAppUser(appUser *user.AppUser) {
	dto.Id = appUser.Id
	dto.UserName = appUser.UserName
	dto.Password = appUser.Password
	dto.Firstname = appUser.Firstname
	dto.Lastname = appUser.Lastname
	dto.Email = appUser.Email
}

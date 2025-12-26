package task

import "apitest/internal/core/user"

type Assignment struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
	TaskId int `json:"task_id"`

	User *user.AppUser
	Task *Task
}

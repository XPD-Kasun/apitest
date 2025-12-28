package gql

import (
	"apitest/internal/core/task"
	"apitest/internal/core/user"
	"strconv"
	"time"
)

func FromAppUser(appuser *user.AppUser) *User {
	return &User{
		ID:   strconv.Itoa(appuser.Id),
		Name: appuser.UserName,
		Age:  20,
	}
}

func ToAppUser(u *User) *user.AppUser {
	id, err := strconv.Atoi(u.ID)
	if err != nil {
		panic(err)
	}
	return &user.AppUser{
		Id:       id,
		UserName: u.Name,
	}
}

func ToTask(taskInput *TaskCreateInput) *task.Task {
	return &task.Task{
		Taskname:    taskInput.Name,
		Description: taskInput.Descript,
		CreatedAt:   time.Time(taskInput.CreatedAt),
		Due:         time.Time(taskInput.Due),
	}
}

func ToGraphTask(t *task.Task) *Task {
	return &Task{
		ID:        strconv.Itoa(t.Id),
		Name:      t.Taskname,
		Descript:  t.Description,
		CreatedAt: Date(t.CreatedAt),
		Due:       Date(t.Due),
		Done:      t.Done,
	}
}

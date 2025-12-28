package sqldto

import (
	"apitest/internal/core/task"

	"github.com/uptrace/bun"
)

type AssignmentSqlDto struct {
	bun.BaseModel `bun:"table:assignment"`

	Id     int `bun:"id"`
	UserId int `bun:"user_id"`
	TaskId int `bun:"task_id"`
}

func (a *AssignmentSqlDto) ToDomainAssignment() *task.Assignment {
	return &task.Assignment{
		Id:     a.Id,
		UserId: a.UserId,
		TaskId: a.TaskId,
	}
}

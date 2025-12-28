package sqldto

import (
	"apitest/internal/core/task"
	"time"

	"github.com/uptrace/bun"
)

type TaskSqlDto struct {
	bun.BaseModel `bun:"table:task"`
	Id            int
	Taskname      string
	Description   string
	CreatedAt     time.Time
	Due           time.Time
	Done          bool
}

func (t *TaskSqlDto) ToCoreTask() *task.Task {
	return &task.Task{
		Id:          t.Id,
		Taskname:    t.Taskname,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		Due:         t.Due,
		Done:        t.Done,
	}
}

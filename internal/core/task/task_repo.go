package task

import (
	"apitest/internal/core/common/baserepo"
)

type TaskRepo interface {
	baserepo.SingleIdGetter[*Task, int]
	baserepo.PaginatedGetter[*Task, int]
	baserepo.Inserter[*Task]
	AssignTask(taskId, userId int) (*Assignment, error)
	GetTasksForUser(id int) ([]*Task, error)
	GetAssignments(taskId int) ([]*Assignment, error)
}

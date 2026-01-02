package task

import (
	"apitest/internal/core/common/baserepo"
)

type TaskService interface {
	CreateNew(task *Task) (*Task, error)
	Assign(taskId, userId int) (*Assignment, error)
	GetTasksForUser(userId int) ([]*Task, error)
	RemoveTaskFromUser(userId, taskId int) error
	GetTasks(cursor, limit int) (*baserepo.PaginatedResult[*Task, int], error)
	GetAssignments(taskId int) ([]*Assignment, error)
	GetTasksByIds(ids []int) ([]*Task, error)
	GetTaskById(id int) (*Task, error)
}

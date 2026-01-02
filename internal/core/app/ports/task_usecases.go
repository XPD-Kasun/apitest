package ports

import (
	"apitest/internal/core/common/baserepo"
	"apitest/internal/core/task"
)

type TaskUseCase interface {
	CreateNewTask(task *task.Task) (*task.Task, error)
	GetTasks(cursor int, limit int) (*baserepo.PaginatedResult[*task.Task, int], error)
	GetTasksForUser(userId int) ([]*task.Task, error)
	GetAssignmentsForTask(taskId int) ([]*task.Assignment, error)
	AssignTaskToUser(taskId, userId int) (*task.Assignment, error)
	GetTasksByIds(ids []int) ([]*task.Task, error)
	GetTaskById(id int) (*task.Task, error)
}

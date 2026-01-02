package usecases

import (
	"apitest/internal/core/app/ports"
	"apitest/internal/core/common/baserepo"
	"apitest/internal/core/task"
)

func NewTaskUseCase(taskSvc task.TaskService) *taskUseCase {
	return &taskUseCase{
		taskSvc: taskSvc,
	}
}

var s ports.TaskUseCase = &taskUseCase{}

type taskUseCase struct {
	taskSvc task.TaskService
}

// GetTaskById implements [ports.TaskUseCase].
func (t *taskUseCase) GetTaskById(id int) (*task.Task, error) {
	return t.taskSvc.GetTaskById(id)
}

// GetTasksByIds implements [ports.TaskUseCase].
func (t *taskUseCase) GetTasksByIds(ids []int) ([]*task.Task, error) {
	return t.taskSvc.GetTasksByIds(ids)
}

// AssignTaskToUser implements [ports.TaskUseCase].
func (t *taskUseCase) AssignTaskToUser(taskId int, userId int) (*task.Assignment, error) {
	return t.taskSvc.Assign(taskId, userId)
}

// GetAssignmentsForTask implements [ports.TaskUseCase].
func (t *taskUseCase) GetAssignmentsForTask(taskId int) ([]*task.Assignment, error) {
	return t.taskSvc.GetAssignments(taskId)
}

// GetTasks implements [ports.TaskUseCase].
func (t *taskUseCase) GetTasks(cursor int, limit int) (*baserepo.PaginatedResult[*task.Task, int], error) {
	return t.taskSvc.GetTasks(cursor, limit)
}

// CreateNewTask implements [ports.TaskUseCase].
func (t *taskUseCase) CreateNewTask(task *task.Task) (*task.Task, error) {
	return t.taskSvc.CreateNew(task)
}

// Assign implements [task.TaskService].
func (t *taskUseCase) Assign(taskId int, userId int) (*task.Assignment, error) {
	return t.taskSvc.Assign(taskId, userId)
}

// CreateNew implements [task.TaskService].
func (t *taskUseCase) CreateNew(task *task.Task) (*task.Task, error) {
	return t.taskSvc.CreateNew(task)
}

// GetTasksForUser implements [task.TaskService].
func (t *taskUseCase) GetTasksForUser(userId int) ([]*task.Task, error) {
	return t.taskSvc.GetTasksForUser(userId)
}

// RemoveTaskFromUser implements [task.TaskService].
func (t *taskUseCase) RemoveTaskFromUser(userId int, taskId int) error {
	panic("unimplemented")
}

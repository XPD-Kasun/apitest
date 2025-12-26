package task

import "apitest/internal/core/common/baserepo"

var a TaskService = &TaskServiceImpl{}

func NewTaskServiceImpl(taskRepo TaskRepo) *TaskServiceImpl {
	return &TaskServiceImpl{
		repo: taskRepo,
	}
}

type TaskServiceImpl struct {
	repo TaskRepo
}

// GetTasks implements [TaskService].
func (t *TaskServiceImpl) GetTasks(cursor int, limit int) (*baserepo.PaginatedResult[*Task, int], error) {
	return t.repo.GetByPage(&baserepo.PaginatedFilter[int]{Cursor: cursor, Limit: limit})
}

// GetAssignments implements [TaskService].
func (t *TaskServiceImpl) GetAssignments(taskId int) ([]*Assignment, error) {
	return t.repo.GetAssignments(taskId)
}

// Assign implements TaskService.
func (t *TaskServiceImpl) Assign(taskId int, userId int) (*Assignment, error) {
	return t.repo.AssignTask(taskId, userId)
}

// CreateNew implements TaskService.
func (t *TaskServiceImpl) CreateNew(task *Task) (*Task, error) {
	err := t.repo.Insert(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// GetTasksForUser implements TaskService.
func (t *TaskServiceImpl) GetTasksForUser(userId int) ([]*Task, error) {
	return t.repo.GetTasksForUser(userId)
}

// RemoveTaskFromUser implements TaskService.
func (t *TaskServiceImpl) RemoveTaskFromUser(userId int, taskId int) error {
	panic("unimplemented")
}

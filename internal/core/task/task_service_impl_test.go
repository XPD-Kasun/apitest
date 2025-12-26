package task

import (
	"apitest/internal/core/common/baserepo"
	"errors"
	"testing"
	"time"
)

var task1 = Task{
	Id:          1,
	Taskname:    "Write thesis draft",
	Description: "Prepare initial draft for the ML hybrid model chapter",
	CreatedAt:   time.Date(2025, 12, 10, 9, 0, 0, 0, time.UTC),
	Due:         time.Date(2025, 12, 20, 23, 59, 0, 0, time.UTC),
}

var task2 = Task{
	Id:          2,
	Taskname:    "Refactor Go filters",
	Description: "Clean up visitor-based NOT pushdown logic",
	CreatedAt:   time.Date(2025, 12, 12, 14, 30, 0, 0, time.UTC),
	Due:         time.Date(2025, 12, 18, 18, 0, 0, 0, time.UTC),
}

var s TaskRepo = &taskRepoImpl{}

type taskRepoImpl struct {
	CanAssign      bool
	CanGetAll      bool
	CanGetUserTask bool
	CanInsert      bool
}

// GetAssignments implements [TaskRepo].
func (t *taskRepoImpl) GetAssignments(taskId int) ([]*Assignment, error) {
	panic("unimplemented")
}

// GetById implements [TaskRepo].
func (t *taskRepoImpl) GetById(id int) (*Task, error) {
	panic("unimplemented")
}

// GetByPage implements [TaskRepo].
func (t *taskRepoImpl) GetByPage(filter *baserepo.PaginatedFilter[int]) (*baserepo.PaginatedResult[*Task, int], error) {
	panic("unimplemented")
}

// AssignTask implements TaskRepo.
func (t *taskRepoImpl) AssignTask(taskId int, userId int) (*Assignment, error) {
	if t.CanAssign {
		return &Assignment{Id: 1, TaskId: taskId, UserId: userId}, nil
	} else {
		return nil, errors.New("cannot assign")
	}

}

// GetAllTasks implements TaskRepo.
func (t *taskRepoImpl) GetAllTasks() ([]*Task, error) {
	if t.CanGetAll {
		return []*Task{&task1, &task2}, nil
	} else {
		return nil, errors.New("error occured while getting dat")
	}
}

// GetTasksForUser implements TaskRepo.
func (t taskRepoImpl) GetTasksForUser(id int) ([]*Task, error) {
	if t.CanGetUserTask {
		return []*Task{&task1, &task2}, nil
	} else {
		return nil, errors.New("error occured while getting data")
	}
}

// Insert implements TaskRepo.
func (t taskRepoImpl) Insert(val *Task) error {
	if t.CanInsert {
		val.Id = 300
		return nil
	} else {
		val.Id = 120
	}
	return errors.New("cannot insert")
}

func TestTaskServiceImpl_createnew(t *testing.T) {

	repo := &taskRepoImpl{}

	taskSvc := TaskServiceImpl{
		repo: repo,
	}

	t.Run("creates new task", func(t *testing.T) {
		repo.CanInsert = true
		task, err := taskSvc.CreateNew(&task1)
		if err != nil {
			t.Error("Cannot create new task")
		}
		if task.Id != 300 {
			t.Error("doesn't set returned id for task")
		}
	})

	t.Run("doesnt create new task", func(t *testing.T) {
		repo.CanInsert = false

		task, err := taskSvc.CreateNew(&task1)
		if err == nil {
			t.Error("no error raised when create new task failed")
		}
		if task != nil {
			t.Error("cannot create task when create new task failed")
		}
	})

}

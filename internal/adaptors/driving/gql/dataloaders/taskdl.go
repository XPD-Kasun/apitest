package dataloaders

import (
	"apitest/internal/core/app/ports"
	"apitest/internal/core/task"
	"context"
)

type TaskDataloader struct {
	TaskUC ports.TaskUseCase
}

func (t *TaskDataloader) GetTasks(ctx context.Context, ids []int) ([]*task.Task, []error) {
	if ctx.Err() != nil {
		return nil, cloneError(ctx.Err(), len(ids))
	}
	tasks, err := t.TaskUC.GetTasksByIds(ids)
	return tasks, cloneError(err, len(tasks))
}

func (t *TaskDataloader) GetTask(ctx context.Context, id int) (*task.Task, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return t.TaskUC.GetTaskById(id)
}

func (t *TaskDataloader) GetAssignments(ctx context.Context, ids []int) ([][]*task.Assignment, []error) {
	if ctx.Err() != nil {
		return nil, cloneError(ctx.Err(), len(ids))
	}
	e := make([]error, len(ids))
	tasks := make([][]*task.Assignment, len(ids))

	for i, id := range ids {
		rtasks, err := t.TaskUC.GetAssignmentsForTask(id)
		e[i] = err
		tasks[i] = rtasks
	}

	return tasks, e
}

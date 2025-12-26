package sql

import (
	"apitest/internal/core/common/baserepo"
	"apitest/internal/core/task"
	"database/sql"
	"errors"
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
)

var (
	NoSingleTaskErr = errors.New("No single task found")
)

func NewMySqlTaskRepo(db *sql.DB) (*MySqlTaskRepo, error) {
	return &MySqlTaskRepo{db: db}, nil
}

var s task.TaskRepo = &MySqlTaskRepo{}

type MySqlTaskRepo struct {
	db *sql.DB
}

// GetAssignments implements [task.TaskRepo].
func (t *MySqlTaskRepo) GetAssignments(taskId int) ([]*task.Assignment, error) {
	cols := []string{"id", "user_id", "task_id"}
	q := sq.Select(cols...).Where(sq.Eq{"task_id": taskId})
	sqlStr, args := q.MustSql()
	log.Debug().Str("sql", sqlStr).Interface("args", args).Msg("GetAssignment")

	rows, err := q.RunWith(t.db).Query()
	if err != nil {
		return nil, err
	}
	return scanAssignments(rows, cols)
}

// GetById implements [task.TaskRepo].
func (t *MySqlTaskRepo) GetById(id int) (*task.Task, error) {

	cols := []string{"id", "taskname", "description", "created_at", "due", "done"}
	rows, err := sq.Select(cols...).From("task").RunWith(t.db).Query()
	if err != nil {
		return nil, err
	}

	tasks, err := scanTasks(rows, cols)
	if len(tasks) != 1 {
		log.Error().Err(err).Msg("at scan tasks")
		return nil, NoSingleTaskErr
	}

	return tasks[0], err

}

// GetByPage implements [task.TaskRepo].
func (t *MySqlTaskRepo) GetByPage(filter *baserepo.PaginatedFilter[int]) (*baserepo.PaginatedResult[*task.Task, int], error) {

	cols := []string{"id", "taskname", "description", "created_at", "due", "done"}
	rows, err := sq.Select(cols...).From("task").Where(sq.GtOrEq{"id": filter.Cursor}).OrderBy("id asc").
		Limit(uint64(filter.Limit)).RunWith(t.db).Query()
	if err != nil {
		return nil, err
	}
	tasks, err := scanTasks(rows, cols)
	result := baserepo.PaginatedResult[*task.Task, int]{
		Items:      tasks,
		HasMore:    true,
		NextCursor: filter.Limit + filter.Cursor,
	}

	return &result, err
}

// AssignTask implements TaskRepo.
func (t *MySqlTaskRepo) AssignTask(taskId int, userId int) (*task.Assignment, error) {
	res, err := sq.Insert("assignment").Columns("taskid", "userid").Values(taskId, userId).
		RunWith(t.db).Exec()
	if err != nil {
		log.Err(err).Msg("cannot insert new assignment into db.")
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Err(err).Msg("cannot retrieve the last inserted id for AssignTask.")
	}
	var assignment = &task.Assignment{
		Id:     int(id),
		TaskId: taskId,
		UserId: userId,
	}

	return assignment, nil
}

// GetAllTasks implements TaskRepo.
func (t *MySqlTaskRepo) GetAllTasks() ([]*task.Task, error) {

	cols := []string{"id", "taskname", "description", "created_at", "due", "done"}
	rows, err := sq.Select(cols...).From("task").RunWith(t.db).Query()
	if err != nil {
		return nil, err
	}
	return scanTasks(rows, cols)
}

// GetTasksForUser implements TaskRepo.
func (t MySqlTaskRepo) GetTasksForUser(userId int) ([]*task.Task, error) {
	cols := []string{"id", "created_at", "taskname", "description", "due", "done"}
	rows, err := sq.Select().
		From("assignment a").
		Join("inner join users u on a.userid = u.id").
		Where("u.id == ?", userId).RunWith(t.db).Query()

	if err != nil {
		log.Err(err).Msg("error running sql")
		return nil, err
	}

	return scanTasks(rows, cols)
}

// Insert implements TaskRepo.
func (t MySqlTaskRepo) Insert(val *task.Task) error {

	result, err := sq.Insert("task").Columns("taskname", "description", "created_at", "due").
		Values(val.Taskname, val.Description, val.CreatedAt, val.Due).RunWith(t.db).Exec()

	if err != nil {
		log.Err(err).Msg("inserting new task failed.")
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Err(err).Msg("getting the inserted task new id failed")
	}

	val.Id = int(id)

	return nil

}

func scanTasks(rows *sql.Rows, cols []string) ([]*task.Task, error) {

	tasks := make([]*task.Task, 0, 10)

	for rows.Next() {

		var task task.Task
		reflectTask := reflect.ValueOf(task)
		fieldPointers := make([]any, len(cols))

		for _, c := range cols {
			ptr := reflectTask.FieldByName(c).Addr().Pointer()
			fieldPointers = append(fieldPointers, ptr)
		}
		rows.Scan(fieldPointers...)

		tasks = append(tasks, &task)
	}

	return tasks, rows.Err()
}

func scanAssignments(rows *sql.Rows, cols []string) ([]*task.Assignment, error) {

	assignments := make([]*task.Assignment, 0, 10)

	for rows.Next() {

		var assign task.Assignment
		reflectTask := reflect.ValueOf(assign)
		fieldPointers := make([]any, len(cols))

		for _, c := range cols {
			ptr := reflectTask.FieldByName(c).Addr().Pointer()
			fieldPointers = append(fieldPointers, ptr)
		}
		rows.Scan(fieldPointers...)

		assignments = append(assignments, &assign)
	}

	return assignments, rows.Err()
}

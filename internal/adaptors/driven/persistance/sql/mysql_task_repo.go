package sql

import (
	"apitest/internal/adaptors/driven/persistance/sql/sqldto"
	"apitest/internal/core/common/baserepo"
	"apitest/internal/core/common/funcs"
	"apitest/internal/core/task"
	"apitest/internal/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

var (
	NoSingleTaskErr = errors.New("No single task found")
)

func NewMySqlTaskRepo(db *bun.DB) (*MySqlTaskRepo, error) {
	db.AddQueryHook(&QueryHook{})
	return &MySqlTaskRepo{db: db}, nil
}

var s task.TaskRepo = &MySqlTaskRepo{}

type MySqlTaskRepo struct {
	db *bun.DB
}

// GetByIds implements [task.TaskRepo].
func (t *MySqlTaskRepo) GetByIds(ids ...int) ([]*task.Task, error) {
	var tasks []*task.Task
	err := bun.NewSelectQuery(t.db).Where("id in (?)", bun.In(ids)).Model(&tasks).Scan(context.TODO())
	if err != nil {
		logger.Error().Err(err).Msg("Repo.GetByIds failed")
		return nil, err
	}
	return tasks, nil
}

// GetAssignments implements [task.TaskRepo].
func (t *MySqlTaskRepo) GetAssignments(taskId int) ([]*task.Assignment, error) {
	var assigns []*sqldto.AssignmentSqlDto
	err := t.db.NewSelect().Model(&assigns).Where("task_id = ?", taskId).Scan(context.TODO())
	if err != nil {
		logger.Error().Err(err).Msg("MySqlTaskRepo::GetAssignments had error")
		return nil, err
	}
	return funcs.Map(assigns, (*sqldto.AssignmentSqlDto).ToDomainAssignment), nil
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
	var tasks []*sqldto.TaskSqlDto
	err := t.db.NewSelect().Model(&tasks).
		Where("id >= ?", filter.Cursor).
		Limit(filter.Limit).
		Scan(context.TODO())

	if err != nil {
		logger.Error().Err(err).Msg("MySqlTaskRepo::GetByPage")
		return nil, err
	}

	for _, t := range tasks {
		fmt.Println(t)
	}

	result := baserepo.PaginatedResult[*task.Task, int]{
		Items:      funcs.Map(tasks, (*sqldto.TaskSqlDto).ToCoreTask),
		HasMore:    true,
		NextCursor: filter.Limit + filter.Cursor,
	}

	return &result, err
}

// AssignTask implements TaskRepo.
func (t *MySqlTaskRepo) AssignTask(taskId int, userId int) (*task.Assignment, error) {

	var totalCols = 0

	r := sq.Select("count(*)").From("assignment	").Where(
		sq.And{
			sq.Eq{"task_id": taskId}, sq.Eq{"user_id": userId},
		}).
		RunWith(t.db).QueryRow()

	err := r.Scan(&totalCols)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get the existing assignments"), err)
	}

	if totalCols > 0 {
		return nil, errors.New("Assignment already exists")
	}

	res, err := sq.Insert("assignment").Columns("task_id", "user_id").Values(taskId, userId).
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

	tasks := make([]*sqldto.TaskSqlDto, 0)
	err := bun.NewRawQuery(t.db, `
		SELECT t.* FROM task t inner join assignment a on t.id = a.task_id 
		where a.user_id = ?
	`, userId).Scan(context.TODO(), &tasks)

	return funcs.Map(tasks, (*sqldto.TaskSqlDto).ToCoreTask), err

	// cols := []string{"id", "created_at", "taskname", "description", "due", "done"}
	// rows, err := sq.Select().
	// 	From("assignment a").
	// 	Join("inner join users u on a.userid = u.id").
	// 	Where("u.id == ?", userId).RunWith(t.db).Query()

	// if err != nil {
	// 	log.Err(err).Msg("error running sql")
	// 	return nil, err
	// }

	// return scanTasks(rows, cols)
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

package common

import (
	"apitest/internal/core/task"
	"apitest/internal/core/user"

	"github.com/vikstrous/dataloadgen"
)

type Key string
type Uniqueid string

var AppContextKey Key = "AppContextKey"

type AppRequestContext struct {
	CorrelationId    Uniqueid
	UserLoader       *dataloadgen.Loader[int, *user.AppUser]
	TaskLoader       *dataloadgen.Loader[int, *task.Task]
	AssignmentLoader *dataloadgen.Loader[int, []*task.Assignment]
}

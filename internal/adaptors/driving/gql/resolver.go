package gql

import (
	"apitest/internal/core/app/ports"
	"apitest/internal/core/common/baserepo"
	"fmt"
	"io"
	"time"
)

type Resolver struct {
	UserUseCase ports.UserUseCase
	TaskUseCase ports.TaskUseCase
}

func NewResolver(userusecase ports.UserUseCase, taskUseCase ports.TaskUseCase) *Resolver {

	return &Resolver{
		UserUseCase: userusecase,
		TaskUseCase: taskUseCase,
	}

}

type PaginatedTasks = baserepo.PaginatedResult[*Task, int]
type Date time.Time

func (date Date) MarshalGQL(w io.Writer) {
	d := time.Time(date)
	datestr := d.Format(time.DateOnly)
	w.Write([]byte(fmt.Sprintf(`"%s"`, datestr)))
}

func (date *Date) UnmarshalGQL(value any) error {
	fmt.Printf("%v %T", value, value)
	dateStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("value %v must be a string", value)
	}
	d, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return err
	}
	*date = Date(d)
	return nil
}

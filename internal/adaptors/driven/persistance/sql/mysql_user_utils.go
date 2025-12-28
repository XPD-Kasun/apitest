package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type QueryHook struct{}

func (h *QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	fmt.Println(time.Since(event.StartTime), string(event.Query))
}

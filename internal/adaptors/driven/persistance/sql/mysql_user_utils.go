package sql

import (
	"apitest/internal/logger"
	"context"

	"github.com/uptrace/bun"
)

type QueryHook struct{}

func (h *QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	logger.Trace().Str("sqlQuery", event.Query).Interface("stashMap", event.Stash).
		Send()
}

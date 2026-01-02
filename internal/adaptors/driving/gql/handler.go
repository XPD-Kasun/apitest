package gql

import (
	"apitest/internal/adaptors/driving/gql/dataloaders"
	"apitest/internal/core/app/ports"
	"apitest/internal/core/common"
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vikstrous/dataloadgen"
)

type loaderSources struct {
	uc ports.UserUseCase
	tc ports.TaskUseCase
}

type ctxKey string

const loaderKey ctxKey = "mxkey"
const loggerKey ctxKey = "loggerkey"

func NewHandler(uc ports.UserUseCase, tu ports.TaskUseCase) gin.HandlerFunc {

	// var sources = loaderSources{
	// 	uc: uc,
	// 	tc: tu,
	// }

	schema := NewExecutableSchema(Config{Resolvers: NewResolver(uc, tu)})
	fmt.Println(schema)
	h := handler.New(schema)
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return func(c *gin.Context) {

		h.ServeHTTP(c.Writer, c.Request)
	}

}

// func GetLoadersFromCtx() lo

func DataLoaderMiddleware(uc ports.UserUseCase, tu ports.TaskUseCase) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		taskdl := dataloaders.TaskDataloader{TaskUC: tu}
		userdl := dataloaders.UserDataloader{UserUC: uc}

		tdl := dataloadgen.NewLoader(taskdl.GetTasks)
		udl := dataloadgen.NewLoader(userdl.GetUsers)
		adl := dataloadgen.NewLoader(taskdl.GetAssignments)
		val := ctx.Request.Context().Value(common.AppContextKey)
		if val != nil {
			if appCtx, ok := val.(*common.AppRequestContext); ok {
				appCtx.TaskLoader = tdl
				appCtx.UserLoader = udl
				appCtx.AssignmentLoader = adl
			}
		}
	}
}

func getAppCtx(ctx context.Context) *common.AppRequestContext {
	val := ctx.Value(common.AppContextKey)
	if val != nil {
		if appCtx, ok := val.(*common.AppRequestContext); ok {
			return appCtx
		}
	}
	panic("getAppCtx retuned error at schemaResolver.Tasks")
}

// Defining the Playground handler
func PlaygroundHandler(gqlUrl string) gin.HandlerFunc {
	h := playground.Handler("GraphQL", gqlUrl)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// func Middleware() gin.HandlerFunc {

// 	return func(c *gin.Context) {

// 		dataloader := dataloadgen.NewLoader()

// 	}
// }

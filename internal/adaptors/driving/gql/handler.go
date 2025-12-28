package gql

import (
	"apitest/internal/core/app/ports"
	"fmt"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

type loaderSources struct {
	uc ports.UserUseCase
	tc ports.TaskUseCase
}

type ctxKey string

const loaderKey ctxKey = "mxkey"

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
		//dataloader := dataloadgen.NewLoader()
		//ctx := context.WithValue(c.Request.Context(), loaderKey)
		//c.Request = c.Request.WithContext(ctx)
		h.ServeHTTP(c.Writer, c.Request)
	}

}

// func GetLoadersFromCtx() lo

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

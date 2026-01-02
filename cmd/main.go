package main

import (
	"apitest/internal/adaptors/driven/persistance/sql"
	"apitest/internal/adaptors/driving/gql"
	"apitest/internal/adaptors/driving/restapi"
	"apitest/internal/adaptors/driving/restapi/middleware"
	"apitest/internal/core/common"
	f "apitest/internal/core/common/filters"
	"apitest/internal/logger"
	"apitest/internal/wiring"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func main2() {
	f := f.OR(
		f.EQ("Id", 5),
		f.AND(
			f.GT("Price", 4000),
			f.BETWEEN("Date", "2024-10-20", time.Now()),
		),
	)

	v := sql.NewSqlVisitor()
	err := f.Accept(v)
	if err != nil {
		fmt.Println(err)
		return
	}

	println(v.String())

}

type SetupServerParams struct {
	fx.In
	Lc           fx.Lifecycle
	Config       *wiring.AppConfig
	AuthCtrl     *restapi.AuthController
	UserCtrl     *restapi.UserController
	TaskCtrl     *restapi.TaskController
	GqlHandler   gin.HandlerFunc `name:"gqlMain"`
	DLMiddleware gin.HandlerFunc `name:"dlMiddleware"`
}

func main() {

	profile, isFound := os.LookupEnv("PROFILE")
	if !isFound {
		fmt.Println("PROFILE is not provided. Setting to dev")
		os.Setenv("PROFILE", "dev")
		profile = "dev"
	}

	if profile != "production" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(".env file not found. Configuration must be available at environment vals.")
		}
	}

	var config wiring.AppConfig
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}

	logger.InitLogger(config.LogLevel)
	fxAppModule := wiring.WireApp(&config)

	fx.New(
		fxAppModule,
		fx.Invoke(setupServer),
	).Run()
}

func setupServer(params SetupServerParams) error {

	url := "localhost:" + strconv.Itoa(int(params.Config.ServerPort))
	engine := gin.New()

	svr := http.Server{
		Addr:    url,
		Handler: engine,
	}

	engine.Use(contextMiddleware())

	authRoute := engine.Group("/api/v1/auth")
	{
		authRoute.Use(middleware.BearerAuthMiddleware())
		authRoute.POST("/login", params.AuthCtrl.Login)
		authRoute.GET("/", params.AuthCtrl.Fn)

	}

	userRoute := engine.Group("/api/v1/users")
	{
		userRoute.POST("/", params.UserCtrl.CreateUser)
	}

	taskRoute := engine.Group("/api/v1/tasks")
	{
		taskRoute.GET("/{id}", params.TaskCtrl.Fn)

	}

	graphQLRoute := engine.Group("/ql/v1")
	{
		graphQLRoute.Use(params.DLMiddleware)
		graphQLRoute.POST("/graphql", params.GqlHandler)
		graphQLRoute.GET("/api", gql.PlaygroundHandler("/ql/v1/graphql"))
	}

	params.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			go func() {
				log.Info().Str("url", url).Msg("server started at given url")
				err := http.ListenAndServe(url, engine)
				if err != nil {
					log.Fatal().Err(err).Msg("http server could not start")
				}
			}()

			return nil

		},
		OnStop: func(ctx context.Context) error {
			return svr.Shutdown(ctx)
		},
	})

	return nil
}

func contextMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		var reqCtx = common.AppRequestContext{}
		setCorrelationId(&reqCtx, c)
		ctx := context.WithValue(c.Request.Context(), common.AppContextKey, &reqCtx)
		c.Request = c.Request.WithContext(ctx)
	}
}

func setCorrelationId(reqCtx *common.AppRequestContext, c *gin.Context) {
	cid := c.Request.Header.Get("X-Correlation-ID")
	if cid == "" {
		cid, err := gonanoid.New()
		if err != nil {
			logger.Error().Err(err).Msg("gonanoid failed to create a new nano id")
		}
		logger.Info().Str("corId", cid).Str("remoteIp", c.Request.RemoteAddr).Msg("BindCorId")
	}
	reqCtx.CorrelationId = common.Uniqueid(cid)
}

func main43() {

	// log.Logger = zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
	// 	w.TimeFormat = time.RFC822
	// 	w.FormatLevel = func(i interface{}) string {
	// 		return strings.ToLower(fmt.Sprintf("%-6s sd", i))
	// 	}
	// 	w.FormatTimestamp = func(i interface{}) string {
	// 		fmt.Println(i, "df")
	// 		return "sdf"
	// 	}

	// }))

	log.Logger = zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.DateTime
		w.NoColor = false

	})).With().Timestamp().Logger()

	log.Error().Msg("sjdfo")

}

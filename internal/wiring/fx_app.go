package wiring

import (
	"apitest/internal/adaptors/driven/persistance/sql"
	"apitest/internal/adaptors/driving/gql"
	"apitest/internal/adaptors/driving/restapi"
	"apitest/internal/core/app/ports"
	"apitest/internal/core/app/usecases"
	"apitest/internal/core/task"
	"apitest/internal/core/user"
	"fmt"

	sql2 "database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"go.uber.org/fx"
)

func WireApp(config *AppConfig) fx.Option {
	var module fx.Option = fx.Module(

		"app",

		fx.Supply(config),

		fx.Provide(
			fx.Annotate(func(config *AppConfig) (*sql2.DB, error) {
				dbConString := fmt.Sprintf("%s:%s@tcp(localhost)/apitest?charset=utf8mb4&parseTime=True&loc=Local", config.DbUserName, config.DbPassword)
				return sql2.Open(config.Provider, dbConString)
			}, fx.ResultTags(`name:"dbConn"`)),
		),

		fx.Provide(
			fx.Annotate(func(db *sql2.DB) *bun.DB {
				return bun.NewDB(db, mysqldialect.New())
			}, fx.ParamTags(`name:"dbConn"`)),
		),

		fx.Provide(
			fx.Annotate(sql.NewMySqlUserRepo, fx.As(new(user.AppUserRepo)), fx.ParamTags(`name:"dbConn"`)),
			fx.Annotate(sql.NewMySqlTaskRepo, fx.As(new(task.TaskRepo))),

			fx.Annotate(func(config *AppConfig, userRepo user.AppUserRepo) *user.UserServiceImpl {
				var usersvc = user.NewUserServiceImpl(userRepo, config.JwtKeyName)
				return usersvc
			}, fx.As(new(user.UserService))),

			fx.Annotate(task.NewTaskServiceImpl, fx.As(new(task.TaskService))),

			fx.Annotate(usecases.NewUserUseCaseImpl, fx.As(new(ports.UserUseCase))),
			fx.Annotate(usecases.NewAuthUseCaseImpl, fx.As(new(ports.AuthUseCase))),
			fx.Annotate(usecases.NewTaskUseCase, fx.As(new(ports.TaskUseCase))),
			//restapi
			restapi.NewAuthCtrl,
			restapi.NewUserCtrl,
			restapi.NewTaskCtrl,
			//graphql
			gql.NewHandler,
		),
	)
	return module
}

// func Wire2(config *AppConfig, setupServer func(*AppConfig, *restapi.Handler)) *fx.App {

// 	var app *fx.App = fx.New(

// 		fx.Supply(config),

// 		fx.Provide(
// 			fx.Annotate(
// 				func(config *AppConfig) (*sql.MySqlUserRepo, error) {
// 					return sql.NewMySqlUserRepo(config.DbConnName)
// 				},
// 				fx.As(new(user.AppUserRepo)),
// 			),
// 		),

// 		fx.Provide(
// 			fx.Annotate(
// 				func(config *AppConfig) (*sql.MySqlTaskRepo, error) {
// 					return sql.NewMySqlTaskRepo(config.DbConnName)
// 				},
// 				fx.As(new(task.TaskRepo)),
// 			),
// 		),

// 		fx.Provide(
// 			fx.Annotate(
// 				func(userRepo user.AppUserRepo, config *AppConfig) *user.UserServiceImpl {
// 					return user.NewUserServiceImpl(userRepo, config.EnvKeyName)
// 				},
// 				fx.As(new(user.UserService)),
// 			),
// 		),

// 		fx.Provide(
// 			fx.Annotate(
// 				func(userSvc user.UserService) *usecases.UserUseCaseImpl {
// 					return usecases.NewUserUseCaseImpl(userSvc)
// 				},
// 				fx.As(new(ports.UserUseCase)),
// 			),
// 		),

// 		fx.Provide(
// 			fx.Annotate(
// 				func(userSvc user.UserService) *usecases.AuthUseCaseImpl {
// 					return usecases.NewAuthUseCaseImpl(userSvc)
// 				},
// 				fx.As(new(ports.AuthUseCase)),
// 			),
// 		),

// 		fx.Provide(func(uc ports.UserUseCase, ac ports.AuthUseCase) *restapi.Handler {
// 			return restapi.NewHandler(ac, uc)
// 		}),

// 		fx.Invoke(func(handler *restapi.Handler) {
// 			setupServer(config, handler)
// 		}),
// 	)

// 	return app
// }

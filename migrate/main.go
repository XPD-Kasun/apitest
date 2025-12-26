package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

func getArgs() (string, error) {

	dir := flag.String("dir", "", "direction of migration")
	flag.Parse()
	if *dir == "" {
		return "", errors.New("specify direction of migration with -dir")
	}

	return *dir, nil
}

type s interface{}

func main() {

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	dir, err := getArgs()
	log.Error().Err(err).Send()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1)/events", username, password))
	log.Error().Err(err).Send()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	log.Error().Err(err).Send()

	migration, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	log.Error().Err(err).Send()

	version, dirty, err := migration.Version()
	if !errors.Is(err, migrate.ErrNilVersion) {
		log.Error().Err(err).Send()
	}

	if dirty {
		migration.Force(int(version))
	}

	switch dir {
	case "up":
		err = migration.Up()
	case "down":
		err = migration.Down()
	default:
		err = errors.New("direction should be either up or down")
	}
	log.Error().Err(err).Send()
}

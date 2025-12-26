package main

import (
	"database/sql"
	"errors"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	dir, err := getArgs()
	logError(err)

	db, err := sql.Open("mysql", "xpd:XPD@tcp(127.0.0.1)/events")
	logError(err)

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	logError(err)

	migration, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	logError(err)

	version, dirty, err := migration.Version()
	if !errors.Is(err, migrate.ErrNilVersion) {
		logError(err)
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
	logError(err)
}

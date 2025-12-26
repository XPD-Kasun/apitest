package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main45() {
	db, err := sql.Open("mysql", "xpd:XPD@tcp(localhost)/apitest")
	if err != nil {
		fmt.Println(err)
		return
	}
	r, e := db.Exec("alter table task modify column done bit not null")
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(r)
}

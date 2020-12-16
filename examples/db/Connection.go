package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

)

func get() *sql.DB {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "secret"
    dbName := "dbcam"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@"+tcp(127.0.0.1:3306)+"/"+dbName)
    if err := db.PingContext(ctx); err != nil {
      log.Fatal(err)
    }
    return db
}

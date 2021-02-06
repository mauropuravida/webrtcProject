package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

)

func get() *sql.DB {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "secret"
    dbName := "dbcam"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@"+"tcp(127.0.0.1:3306)"+"/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

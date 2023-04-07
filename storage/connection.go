package storage

import (
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
)

func ConnectionToDB() (dbr.Session, error) {

	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1",
		DBName:               "hrms_oopm",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	/* 	var connectionString string

	   	if database == "mariadb" {
	   		connectionString = ("root: @tcp(127.0.0.1)/hrms_oopm")
	   	} else if database == "maria2" {
	   		connectionString = ("root@local:@tcp(127.0.0.1)/hrms_oopm")
	   	} else {
	   		return dbr.Session{}, fmt.Errorf("invalid database name: %s", database)
	   	} */

	conn, err := dbr.Open("mysql", cfg.FormatDSN(), nil)
	if err != nil {
		return dbr.Session{}, err
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	session := conn.NewSession(nil)

	return *session, nil
}

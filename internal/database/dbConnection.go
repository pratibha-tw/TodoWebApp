package database

import (
	"database/sql"
	"fmt"
	"log"
	"todoapp/config"
)

func CreateConnection(cfg config.Config) *sql.DB {
	config.GetConfigs(&cfg)
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)

	fmt.Println(connectionString)
	dbConn, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("unable to open connection with database ", err.Error())
	}
	if err := dbConn.Ping(); err != nil {
		log.Fatal("unable to ping database ", err.Error())
	}
	return dbConn
}

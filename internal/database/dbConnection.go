package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"todoapp/config"

	"github.com/redis/go-redis/v9"
)

func CreateConnection(cfg config.Config) *sql.DB {
	//config.GetConfigs(&cfg)
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Database.User, cfg.Database.Password,
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

func CreateRedisConnection(cfg config.Config) *redis.Client {

	ctx := context.Background()
	//fmt.Println("Test Redis connection")
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ping, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error in connecting redis")
		return nil
	}
	fmt.Printf("successfully connected to redis %s", ping)
	return client
}

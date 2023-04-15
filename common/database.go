package common

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func InitDB() *gorm.DB {
	dbPassword := GetEnv("DB_PASSWORD")
	dbHost := GetEnv("DB_HOST")
	dbPort := GetEnv("DB_PORT")
	dbName := GetEnv("DB_NAME")
	dbUser := GetEnv("DB_USER")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		panic(err)
	}

	postgresdb, err := db.DB()
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		panic(err)
	}

	postgresdb.SetMaxOpenConns(100)

	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
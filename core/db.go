package core

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabase() *gorm.DB {
	dsn := "host=localhost user=dful password=Tacos123.# dbname=horus port=5432 sslmode=disable"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Wrong database Conection")
	}

	sqldb, err := conn.DB()

	err = sqldb.Ping()
	if err != nil {
		log.Fatalln("Database Conection error")
	}

	return conn
}

func InitDB() {
	conn := GetDatabase()
	defer CloseDB(conn)
	conn.AutoMigrate(User{})
}

func CloseDB(conn *gorm.DB) {
	sqldb, _ := conn.DB()
	sqldb.Close()
}

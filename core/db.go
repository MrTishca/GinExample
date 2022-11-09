package core

import (
	"log"

	"gorm.io/driver/postgres"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func GetDatabase() *gorm.DB {
	
	host := viperEnvVariable("Host")
	user := viperEnvVariable("User")
	password := viperEnvVariable("Password")
	dbname := viperEnvVariable("DBname")
	port := viperEnvVariable("Port")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"
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

func viperEnvVariable(key string)string{
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil{
		log.Fatalf("Errorr",err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Not Found Value")
	}
	return value
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

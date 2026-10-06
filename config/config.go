package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}
type Jwt struct {
	ACCESS_TOKEN_SECRET_KEY  string
	REFRESH_TOKEN_SECRET_KEY string
}

type Server struct {
	Host string
	Port string
}

type Config struct {
	Server
	Database
	Jwt
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(err.Error())
		return nil
	}

	return &Config{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: Database{
			Driver:   os.Getenv("DB_DRIVER"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASS"),
			DBName:   os.Getenv("DB_NAME"),
		},
		Jwt: Jwt{
			REFRESH_TOKEN_SECRET_KEY: os.Getenv("JWT_REFRESH_SECRET_KEY"),
			ACCESS_TOKEN_SECRET_KEY:  os.Getenv("JWT_ACCESS_SECRET_KEY"),
		},
	}
}

func (conf *Config) GetConnectionDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DBName,
	)

	db, err := gorm.Open(mysqlDriver.Open(dsn))
	if err != nil {
		log.Println(err.Error())
	}

	return db
}

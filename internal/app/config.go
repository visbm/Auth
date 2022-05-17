package app

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Host string
		Port int
	}
	SQLLite3DB struct {
		Host     string
		Port     int
		Username string
		Password string
		DbName   string
		SSLMODE  string
	}
}

func (с Config) init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Print(err)
	}
}

func (c *Config) NewConfig() {
	c.init()

	c.Server.Host = os.Getenv("SERVER_HOST")
	c.Server.Port, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))

	
	c.SQLLite3DB.Host = os.Getenv("SQLITE3_HOST")
	c.SQLLite3DB.DbName = os.Getenv("SQLITE3_DB")
	c.SQLLite3DB.Username = os.Getenv("SQLITE3_USER")
	c.SQLLite3DB.Password = os.Getenv("SQLITE3_PASSWORD")
	c.SQLLite3DB.Port, _ = strconv.Atoi(os.Getenv("SQLITE3_PORT"))
	c.SQLLite3DB.SSLMODE = os.Getenv("SQLITE3_SSLMODE")
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
func (c *Config) SQLLite3DataSourceName() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.SQLLite3DB.Host,
		c.SQLLite3DB.Port,
		c.SQLLite3DB.Username,
		c.SQLLite3DB.Password,
		c.SQLLite3DB.DbName,
		c.SQLLite3DB.SSLMODE,
	)
}

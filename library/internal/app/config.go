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
		Path     string
		Username string
		Password string
	}
}

func (с Config) init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Print(err)
	}
}

func (c *Config) NewConfig() {
	c.init()

	c.Server.Host = os.Getenv("SERVER_HOST")
	c.Server.Port, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))

	
	c.SQLLite3DB.Path = os.Getenv("SQLITE3_PATH")
	c.SQLLite3DB.Username = os.Getenv("SQLITE3_USER")
	c.SQLLite3DB.Password = os.Getenv("SQLITE3_PASSWORD")
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
func (c *Config) SQLLite3DataSourceName() string {
	return fmt.Sprintf(
		"%s?_auth&_auth_user=%s&_auth_pass=%s",
		c.SQLLite3DB.Path,
		c.SQLLite3DB.Username,
		c.SQLLite3DB.Password,
	)
}

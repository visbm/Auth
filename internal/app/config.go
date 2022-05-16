package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Host     string
		Port     int
	}
}

func(с Config) init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Print(err)
	}
}

func (c *Config) NewConfig() {
	c.init()

	c.Server.Host = os.Getenv("SERVER_HOST")
	c.Server.Port, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Env struct {
	HOST            string
	PORT            string
	POSTGRES        *Postgres
	ALLOWED_ORIGINS []string
	REDIS           *Redis
}

func LoadEnv() *Env {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("[ConfigEnv.LoadEnv] load : %v", err.Error())
	}

	origins := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	for _, o := range strings.Split(origins, ",") {
		trimmed := strings.Trim(o, " \t\n\r\"'")
		if trimmed != "" {
			allowedOrigins = append(allowedOrigins, trimmed)
		}
	}

	return &Env{
		HOST: os.Getenv("HOST_SERVER"),
		PORT: os.Getenv("PORT_SERVER"),
		POSTGRES: &Postgres{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Db_name:  os.Getenv("POSTGRES_DB"),
		},
		REDIS: &Redis{
			HOST: os.Getenv("REDIS_HOST"),
			PORT: os.Getenv("REDIS_PORT"),
			USER: os.Getenv("REDIS_USER"),
			PWD:  os.Getenv("REDIS_PASWORD"),
		},
		ALLOWED_ORIGINS: allowedOrigins,
	}
}

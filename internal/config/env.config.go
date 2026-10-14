package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	HOST     string
	PORT     string
	POSTGRES *Postgres
	REDIS    *Redis
}

func LoadEnv() *Env {
	if err := godotenv.Load(); err != nil {
		log.Printf("[ConfigEnv.LoadEnv] load : %v", err.Error())
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
			PWD:  os.Getenv("REDIS_PASSWORD"),
		},
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/MGT06/e-wallet_backend/internal/config"
	"github.com/MGT06/e-wallet_backend/internal/router"
	"github.com/gin-gonic/gin"
)

// @title           			Application E-Walet
// @version         			1.0
// @description     			This is a sample send request use API.

// @host      					192.168.0.186:5000
// @BasePath  					/

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description					Bearer Token used as identity for accessing backend
func main() {
	// ENV
	conf := config.LoadEnv()

	// Redis
	redis := conf.REDIS.Connect()
	defer func() {
		if err := redis.Close(); err != nil {
			log.Println("Close Redis : ", err.Error())
		}
	}()

	// PostgreSQL
	pool := conf.POSTGRES.Connect()
	defer pool.Close()

	// Gin Gonic
	// gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	router.NewRouter(server, pool, redis, conf.ALLOWED_ORIGINS).Connect()

	server.Run(fmt.Sprintf("%s:%s", conf.HOST, conf.PORT))
}

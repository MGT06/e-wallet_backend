package main

import (
	"fmt"
	"log"
	"time"
	"os"

	"github.com/MGT06/e-wallet_backend/docs"
	"github.com/MGT06/e-wallet_backend/internal/config"
	"github.com/MGT06/e-wallet_backend/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// @title           			Application E-Walet
// @version         			1.0
// @description     			This is a sample send request use API.

// @host      					localhost:5000
// @BasePath  					/

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description					Bearer Token used as identity for accessing backend
func main() {
	// ENV
	conf := config.LoadEnv()
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", os.Getenv("HOST_SWAGG"), os.Getenv("PORT_SWAGG"))

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

	// Resty
	client := resty.New().SetTimeout(30 * time.Second).SetRetryCount(2).SetRetryWaitTime(1 * time.Second)
	resp, err := client.R().
		EnableTrace().
		Get("https://httpbin.org/get")

	if err != nil {
		log.Println(err)
	}

	log.Println("Healty resty client: ", resp.IsSuccess())

	// Gin Gonic
	// gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	router.Router(server, pool, redis, client)

	server.Run(fmt.Sprintf("%s:%s", conf.HOST, conf.PORT))
}

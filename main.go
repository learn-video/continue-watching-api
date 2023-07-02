package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/learn-video/continue-watching-api/position"
	"github.com/redis/go-redis/v9"
)

type config struct {
	RedisHost string `env:"REDIS_HOST"`
	RedisPort int    `env:"REDIS_PORT"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error loading config: %s", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: "",
		DB:       0,
	})

	handler := position.NewHandler(rdb)

	e := echo.New()
	e.Use(middleware.Logger())
	e.POST("/watching", handler.Record)
	e.GET("/watching", handler.Fetch)
	e.Logger.Fatal(e.Start(":8000"))
}

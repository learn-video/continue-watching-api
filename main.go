package main

import (
	"github.com/labstack/echo/v4"
	"github.com/learn-video/continue-watching-api/position"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	handler := position.NewHandler(rdb)

	e := echo.New()
	e.POST("/watching", handler.Record)
	e.Logger.Fatal(e.Start(":8000"))
}

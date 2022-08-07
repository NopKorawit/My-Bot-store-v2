package main

import (
	"store/pkg/config"
	"store/pkg/line"
	"store/pkg/sheet"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.New()

	sheetService, err := sheet.NewService(cfg)
	if err != nil {
		panic(err)
	}

	bot := line.NewBot(sheetService, cfg)

	route := gin.Default()
	route.Use(cors.Default())
	route.POST("/api/v1/line/callback", bot.Callback)
	route.Run()
}

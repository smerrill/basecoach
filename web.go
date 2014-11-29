package main

import (
	"github.com/gin-gonic/gin"
)

func RunServer(config Config, checks map[string]Check) {
	InitHttpClient(config.RequestTimeout)

	router := gin.Default()

	// @TODO: Make a nicer homepage.
	router.GET("/", func(context *gin.Context) {
		context.String(200, "Base Coach is waving you home on %s.", config.BindAddress)
	})

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})

	for _, check := range checks {
		router.GET(check.Url, func(context *gin.Context) {
			check.RunCheck(context, config)
		})
	}

	router.Run(config.BindAddress)
}

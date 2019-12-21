package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
	"io/ioutil"
	"os"
)

func StartServer() error {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	router := gin.Default()
	router.Use(gin.Recovery())
	RegisterHandler(router)
	port:= ":"+os.Getenv("application.port")
	_ = glg.Log("Listening at,", port)
	return router.Run(port)
}

func RegisterHandler(router *gin.Engine){
	router.POST("/", home)
	router.GET("/", home)
}

func home(c *gin.Context) {

	c.JSON(200, gin.H{
		"app.name" : os.Getenv("application.name"),
		"app.desc": os.Getenv("application.desc"),
		"app.ver": os.Getenv("application.ver"),
		"port.listener":os.Getenv("application.port"),
	})
}

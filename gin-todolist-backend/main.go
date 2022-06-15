// main.go
package main

import (
	// "net/http"

	"github.com/gin-gonic/gin"
)
var router *gin.Engine

func main() {

	// Set the router as the default one provided by Gin
	router = gin.Default()
	router.GET("/", renderAllTasks)
	router.GET("/task/:task_id", renderTaskById)
    router.Run()
}
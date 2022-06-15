// routers.go

package main

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	router := gin.Default()
	
	// deprecated
	router.GET("/task/:task_id", renderTaskById)

	// Landing page: render all tasks
	router.GET("/", renderAllTasks)
	
	// Create task
	// router.GET("/create", createTask)

	// Edit task
	// router.GET("/edit", createTask)

	// Delete task
	// router.GET("/delete", createTask)

    return router
}


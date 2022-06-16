// routers.go

package main

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	router := gin.Default()
	
	// deprecated
	// router.GET("/task/:task_id", renderTaskById)

	// Landing page: render all tasks
	router.GET("/", renderAllTasks)
	
	// Create task
	router.POST("/create/:parent_id/:content", createTaskService)

	// Edit task
	// router.PUT("/edit/:id/:new_content", editTaskService)

	// Delete task
	router.DELETE("/delete/:id", deleteTaskService)

    return router
}


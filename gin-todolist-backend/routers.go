// routers.go

package main

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	router := gin.Default()

	// Landing page: Get All Tasks
	router.GET("/tasks", getAllTasksService)

	// Create task
	router.POST("/create/:content", createTaskService)

	// Edit task
	router.PUT("/edit/:id/:new_content", editTaskService)

	// // Delete task
	// router.DELETE("/delete/:id", deleteTaskService)

	// // Create subtask
	// router.POST("/create-subtask/:parent_id/:content", createSubtaskService)

	// // Delete subtask
	// router.DELETE("/delete-subtask/:parent_id/:subtask_id", deleteSubtaskService)

	// // Edit subtask
	// router.PUT("/edit-subtask/:parent_id/:subtask_id/:new_content", editSubtaskService)

	return router
}

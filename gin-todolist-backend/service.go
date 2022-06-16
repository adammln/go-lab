// service.go

package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func _format_handler(c *gin.Context, data []Task) {
	switch c.Request.Header.Get("Accept") {
		case "application/xml":
			c.XML(http.StatusOK, data)
		default:
			c.JSON(http.StatusOK, data)
	}
}

func renderAllTasks(c *gin.Context) {
	tasks := getAllTasks()
	_format_handler(c, tasks)
}

func createTaskService(c *gin.Context) {
	createTask(c.Param("content"))
	renderAllTasks(c)
}

// func editTaskService(c *gin.Context) {
// 	if id, err := strconv.Atoi(c.Param("id")); err == nil {
// 		editTask(id, c.Param("new_content"))
// 	} else {
// 		c.AbortWithError(http.StatusNotFound, err)
// 	}
// 	renderAllTasks(c)
// }

func deleteTaskService(c *gin.Context) {
	deleteTask(c.Param("id"))
	renderAllTasks(c)
}
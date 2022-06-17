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

func editTaskService(c *gin.Context) {
	editTask(c.Param("id"), c.Param("new_content"))
	renderAllTasks(c)
}

func deleteTaskService(c *gin.Context) {
	deleteTask(c.Param("id"))
	renderAllTasks(c)
}

// subtask CRUD
func createSubtaskService(c *gin.Context) {
	createSubtask(c.Param("parent_id"), c.Param("content"))
	renderAllTasks(c)
}

func editSubtaskService(c *gin.Context) {
	editSubtask(c.Param("parent_id"), c.Param("subtask_id"), c.Param("new_content"))
	renderAllTasks(c)
}

func deleteSubtaskService(c *gin.Context) {
	deleteSubtask(c.Param("parent_id"), c.Param("subtask_id"))
	renderAllTasks(c)
}
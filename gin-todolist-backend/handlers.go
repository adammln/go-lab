// handlers.task.go

package main

import (
  "net/http"
  "strconv"
//   "fmt"

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

func renderTaskById(c *gin.Context) {
	if taskId, err := strconv.Atoi(c.Param("task_id")); err == nil {
		task, _ := getTaskById(taskId)
		_format_handler(c, []Task{task})
	} else {
		c.AbortWithError(http.StatusNotFound, err)
	}
}
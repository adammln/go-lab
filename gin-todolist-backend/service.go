// service.go

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

func createTaskService(c *gin.Context) {
	if parentId, err := strconv.Atoi(c.Param("parent_id")); err == nil {
		createTask(parentId, c.Param("content"))
	} else {
		c.AbortWithError(http.StatusNotFound, err)
	}
	renderAllTasks(c)
}

func editTaskService(c *gin.Context) {
	if id, err := strconv.Atoi(c.Param("id")); err == nil {
		editTask(id, c.Param("new_content"))
	} else {
		c.AbortWithError(http.StatusNotFound, err)
	}
	renderAllTasks(c)
}

func deleteTaskService(c *gin.Context) {
	id, errId := strconv.Atoi(c.Param("id"))
	parentId, errParentId := strconv.Atoi(c.Param("parent_id"))

	if (errId == nil && errParentId == nil) {
		deleteTask(id, parentId)
	} else {
		c.AbortWithError(http.StatusNotFound, errId)
	}
	renderAllTasks(c)
}


// deprecated
func renderTaskById(c *gin.Context) {
	if taskId, err := strconv.Atoi(c.Param("task_id")); err == nil {
		task, _ := getTaskById(taskId)
		_format_handler(c, []Task{task})
	} else {
		c.AbortWithError(http.StatusNotFound, err)
	}
}
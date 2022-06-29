// service.go

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var collectionID string = os.Getenv("FIRESTORE_DATA_COLLECTION_ID")

func _response_format_handler(c *gin.Context, data *TaskWrapper) {
	switch c.Request.Header.Get("Accept") {
	case "application/xml":
		c.XML(http.StatusOK, data)
	default:
		c.JSON(http.StatusOK, data)
	}
}

func getAllTasksService(c *gin.Context) {
	tasks, err := dbGetAllTasks(c, collectionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	_response_format_handler(c, tasks)
}

func createTaskService(c *gin.Context) {
	if rank_order, err := strconv.Atoi(c.Param("rank_order")); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		wr, err := dbCreateTask(c, c.Param("content"), rank_order, collectionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, wr)
	}
}

func editTaskService(c *gin.Context) {
	requestBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	wr, err := dbEditTask(c, c.Param("id"), requestBody, collectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, wr)
}

// func deleteTaskService(c *gin.Context) {
// 	deleteTask(c.Param("id"))
// 	renderAllTasks(c)
// }

// // subtask CRUD
// func createSubtaskService(c *gin.Context) {
// 	createSubtask(c.Param("parent_id"), c.Param("content"))
// 	renderAllTasks(c)
// }

// func editSubtaskService(c *gin.Context) {
// 	editSubtask(c.Param("parent_id"), c.Param("subtask_id"), c.Param("new_content"))
// 	renderAllTasks(c)
// }

// func deleteSubtaskService(c *gin.Context) {
// 	deleteSubtask(c.Param("parent_id"), c.Param("subtask_id"))
// 	renderAllTasks(c)
// }

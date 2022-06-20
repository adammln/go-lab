package main

import (
	// "net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func _getTestContext() *gin.Context {
	/** Helper: Generate Gin Framework Test Context
	**/
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	return c
}

func _getDummyTasks() map[string]Task {
	/** Helper: Generate 10 dummy tasks
	- 5 Parent Task (ParentID=nil) 
	- 5 Subtasks
	- 1 Parent -> 3 subtasks
	- 1 Parent -> 2 subtasks
	- 3 Parent -> 0 subtasks (Subtasks=nil)
	**/
	
	// # generate default value
	var tasks map[string]Task
	for i := 1; i <= 10; i++ {
		tmpID := fmt.Sprintf("task_id_%d", i)
		tmpTask := Task{
			ID: tmpID,
			ParentID: nil,
			RankOrder: i,
			Content: fmt.Sprintf("original content for task_id_%d", i),
			IsChecked: false,
			Subtasks: nil,
		}
		tasks[tmpID] = tmpTask
	}

	// # subtasks mapping
	// parent	: subtasks
	// 1			: 2,3,4
	subtasks_1 := []string{"task_id_2", "task_id_3", "task_id_4"}
	tasks["task_id_1"].Subtasks = subtasks_1
	for i, subtaskId := range subtasks_1 {
		tasks[subtaskId].ParentID = "task_id_1"
		tasks[subtaskId].RankOrder = i+1
		fmt.Println(i)
	}

	// 5			: 6,7
	subtasks_5 := []string{"task_id_6", "task_id_7"}
	tasks["task_id_5"].Subtasks = subtasks_5
	for i, subtaskId := range subtasks_1 {
		tasks[subtaskId].ParentID = "task_id_5"
		tasks[subtaskId].RankOrder = i+1
		fmt.Println(i)
	}
	return tasks
}

func TestFirebaseConnectionSuccess(t *testing.T) {
	c := _getTestContext()
	// test firebase app initialization 
	app, errFirebase := _initFirebaseApp(c)
	if (errFirebase != nil) {
		t.Fatal(errFirebase)
	}

	// test firestore client connection
	client, errFirestore := _initFirestoreClient(c, app)
	if (errFirestore != nil) {
		t.Fatal(errFirestore)
	}
	client.Close()
}


func TestGetAllTasks(t *testing.T) {
	c := _getTestContext()
	collection_id := os.Getenv("FIRESTORE_TEST_DATA_COLLECTION_ID")

	WANTED_TASKS_LENGTH := 10

	data := dbGetAllTasks(c)

	// check sample data != nil
	if (data == nil) {
		t.Fatal("data query = NULL, want []Task")
	}

	// TODO: check if length of data as we wanted

	// TODO: check if ordering is correct

	// TODO: check if all attributes are available
}
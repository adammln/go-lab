package main

import (
	// "net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"os"
	"reflect"

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
	var tasks = make(map[string]Task)
	for i := 1; i <= 10; i++ {
		tmpID := fmt.Sprintf("task_id_%d", i)
		tmpTask := Task{
			ID: tmpID,
			ParentID: "",
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
	tmpTask_1 := tasks["task_id_1"]
	tmpTask_1.Subtasks = subtasks_1
	tasks["task_id_1"] = tmpTask_1
	for i, subtaskId := range subtasks_1 {
		tmpSubtask := tasks[subtaskId]
		tmpSubtask.ParentID = "task_id_1"
		tmpSubtask.RankOrder = i+1
		tasks[subtaskId] = tmpSubtask
	}

	// parent	: subtasks
	// 5			: 6,7
	subtasks_5 := []string{"task_id_6", "task_id_7"}
	tmpTask_5 := tasks["task_id_5"]
	tmpTask_5.Subtasks = subtasks_5
	tasks["task_id_1"] = tmpTask_5
	for i, subtaskId := range subtasks_5 {
		tmpSubtask := tasks[subtaskId]
		tmpSubtask.ParentID =  "task_id_5"
		tmpSubtask.RankOrder = i+1
		tasks[subtaskId] = tmpSubtask
	}
	return tasks
}

func TestFirebaseConnection(t *testing.T) {
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

func TestFirebaseConnectionWrapper(t *testing.T) {
	c := _getTestContext()

	// test Firebase (Firestore) connection wrapper 
	db := _initDbConnection(c)
	if db == nil {
		t.Fatal("DB Error: db client is nil")
	}
	db.Close()
}

func TestCreateNewTasks(t *testing.T) {
	c := _getTestContext()
	collectionID := os.Getenv("FIRESTORE_TEST_DATA_COLLECTION_ID")
	tasks := _getDummyTasks()
	for _, task_id := range []string{
		"task_id_1", "task_id_5", "task_id_6", 
		"task_id_7", "task_id_8", "task_id_9",
		"task_id_10",
	} {
		task := tasks[task_id]
		_, err := dbCreateTask(
			c,
			task.Content,
			task.RankOrder,
			collectionID,
		)
		if (err != nil) {
			t.Fatal(err)
		}
	}
}

func TestGetAllTasks(t *testing.T) {
	c := _getTestContext()
	collection_id := os.Getenv("FIRESTORE_TEST_DATA_COLLECTION_ID")

	tasksWrapper := dbGetAllTasks(c)

	// # TEST datatype is correct: 
	// - tasksWrapper == TasksWrapper
	tTaskWrapper := reflect.TypeOf(taskWrapper)
	if (tTaskWrapper != reflect.TypeOf(TaskWrapper)) {
		t.Fatalf(`TypeError: Type of taskWrapper = %s, want %s`, tTaskWrapper, reflect.TypeOf(TaskWrapper))
	}

	// - tasksWrapper.Data -> == map[string]Task
	tTaskWrapperData := reflect.TypeOf(taskWrapper.Data)
	if (tTaskWrapperData != reflect.TypeOf(map[string]Task)) {
		t.Fatalf(`TypeError: Type of taskWrapper.Data = %s, want %s`, tTaskWrapper, reflect.TypeOf(map[string]Task))
	}

	// - tasksWrapper.Orders -> == []string
	tTaskWrapperOrders := reflect.TypeOf(taskWrapper.Orders)
	if (tTaskWrapperOrders != reflect.TypeOf([]string)) {
		t.Fatalf(`TypeError: Type of taskWrapper.Data = %s, want %s`, tTaskWrapperOrders, reflect.TypeOf([]string))
	}
	
	// # TEST query result as wanted
	// - specify wanted length results:
	wantedTasksLength := 5
	actualTasksLengthData := len(tasksWrapper.Data)
	actualTasksLengthOrders := len(tasksWrapper.Orders)
	
	// - check tasksWrapper.Data == 5
	if actualTasksLengthData != WANTED_TASKS_LENGTH {
		t.Fatalf(
			`Result Length of dbGetAllTasks()=>tasksWrapper.Data = %d, want %d, error`, 
			actualTasksLengthData, 
			wantedTasksLength,
		)
	}

	// - check tasksWrapper.Orders == 5
	if actualTasksLengthOrders != WANTED_TASKS_LENGTH {
		t.Fatalf(
			`Result Length of dbGetAllTasks()=>tasksWrapper.Orders = %d, want %d, error`, 
			actualTasksLengthOrders, 
			wantedTasksLength,
		)
	}

	// Iteration check
	i := 0
	for docId, task := range tasksWrapper.Data {
		// # TEST ordering is correct: [1-5]
		if task.RankOrder != i+1 {
			t.Fatalf(`Ordering error: ordering for task "%s" = %d, want %d`, task.Content, task.RankOrder, i+1)
		}

		// # TEST default attr ID: wanted not "" (empty string)
		if task.ID == "" {
			t.Fatalf(`Default value of the attribute 'ID' for task "%s" = "", want %s`, task.Content, task.RankOrder, docId)
		}

		// # TEST default attr IsChecked: wanted false
		if task.IsChecked {
			t.Fatalf(`Default value of the attribute 'IsChecked' for task "%s" = %t, want false`, task.Content, task.IsChecked)
		}
		i++
	}
}

// func TestGetAllTasks(t *testing.T) {
// 	c := _getTestContext()
// 	collection_id := os.Getenv("FIRESTORE_TEST_DATA_COLLECTION_ID")

// 	WANTED_TASKS_LENGTH := 10

// 	data := dbGetAllTasks(c)

// 	// check sample data != nil
// 	if (data == nil) {
// 		t.Fatal("data query = NULL, want []Task")
// 	}

// 	// TODO: check if length of data as we wanted

// 	// TODO: check if ordering is correct

// 	// TODO: check if all attributes are available
// }

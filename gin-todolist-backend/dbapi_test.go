package main

import (
	// "net/http"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	firestore "cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
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
			ID:        tmpID,
			ParentID:  "",
			RankOrder: i,
			Content:   fmt.Sprintf("original content for task_id_%d", i),
			IsChecked: false,
			Subtasks:  nil,
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
		tmpSubtask.RankOrder = i + 1
		tasks[subtaskId] = tmpSubtask
	}

	// parent	: subtasks
	// 5			: 6,7
	subtasks_5 := []string{"task_id_6", "task_id_7"}
	tmpTask_5 := tasks["task_id_5"]
	tmpTask_5.Subtasks = subtasks_5
	tasks["task_id_5"] = tmpTask_5
	for i, subtaskId := range subtasks_5 {
		tmpSubtask := tasks[subtaskId]
		tmpSubtask.ParentID = "task_id_5"
		tmpSubtask.RankOrder = i + 1
		tasks[subtaskId] = tmpSubtask
	}
	return tasks
}

func TestFirebaseConnection(t *testing.T) {
	c := _getTestContext()
	// test firebase app initialization
	app, errFirebase := _initFirebaseApp(c)
	if errFirebase != nil {
		t.Fatal(errFirebase)
	}

	// test firestore client connection
	client, errFirestore := _initFirestoreClient(c, app)
	if errFirestore != nil {
		t.Fatal(errFirestore)
	}
	client.Close()
}

func TestFirebaseConnectionWrapper(t *testing.T) {
	c := _getTestContext()

	// test Firebase (Firestore) connection wrapper
	db := _initDbConnection(c)
	if db == nil {
		t.Fatal("[ERROR] TestFirebaseConnectionWrapper: db client is nil")
	}
	db.Close()
}

func TestCreateNewTasks(t *testing.T) {
	c := _getTestContext()
	collectionID := os.Getenv("FIRESTORE_TEST_DATA_COLLECTION_ID")
	tasks := _getDummyTasks()
	for i, taskID := range []string{
		"task_id_1", "task_id_5", "task_id_8",
		"task_id_9", "task_id_10",
	} {
		task := tasks[taskID]
		_, err := dbCreateTask(
			c,
			task.Content,
			i+1,
			collectionID,
		)
		if err != nil {
			t.Fatalf(`[Error] TestCreateNewTasks: Failed at creating task-%d of %d. Reason: %s`, i, 5, err)
		}
	}
}

func TestGetAllTasks(t *testing.T) {
	c := _getTestContext()
	collectionID := os.Getenv("FIRESTORE_TEST_DATA_COLLECTION_ID")

	taskWrapper, err := dbGetAllTasks(c, collectionID)

	// # TEST err
	if err != nil {
		t.Fatalf(`[ERROR] dbapi_test.TestGetAllTasks: Fail at getting all tasks. %v`, err)
	}

	// # TEST datatype is correct:
	// - taskWrapper == TaskWrapper
	var tw TaskWrapper
	if reflect.TypeOf(*taskWrapper) != reflect.TypeOf(tw) {
		t.Fatalf(`[ERROR] dbapi_test.TestGetAllTasks: Type of taskWrapper = %s, want %s`, reflect.TypeOf(*taskWrapper), reflect.TypeOf(tw))
	}

	// - taskWrapper.Data -> == map[string]Task
	var mpTask map[string]Task
	if reflect.TypeOf(taskWrapper.Data) != reflect.TypeOf(mpTask) {
		t.Fatalf(`[ERROR] dbapi_test.TestGetAllTasks: Type of taskWrapper.Data = %s, want %s`, reflect.TypeOf(taskWrapper.Data), reflect.TypeOf(mpTask))
	}

	// - taskWrapper.Orders -> == []string
	var lstr []string
	if reflect.TypeOf(taskWrapper.Orders) != reflect.TypeOf(lstr) {
		t.Fatalf(`[ERROR] dbapi_test.TestGetAllTasks: Type of taskWrapper.Data = %s, want %s`, reflect.TypeOf(taskWrapper.Orders), reflect.TypeOf(lstr))
	}

	// # TEST query result as wanted
	// - specify wanted length results:
	// actualTasksLengthData := len(taskWrapper.Data)

	// wantedTasksLengthData := 10
	// actualTasksLengthData := len(taskWrapper.Data)
	// // - check taskWrapper.Data == 5 (aborted, because Data includes subtasks)
	// if actualTasksLengthData != wantedTasksLengthData {
	// 	t.Fatalf(
	// 		`Result Length of dbGetAllTasks()=>taskWrapper.Data = %d, want %d, error`,
	// 		actualTasksLengthData,
	// 		wantedTasksLengthData,
	// 	)
	// }

	wantedTasksLengthOrder := 5
	actualTasksLengthOrders := len(taskWrapper.Orders)
	// - check len(taskWrapper.Orders) == 5
	if actualTasksLengthOrders != wantedTasksLengthOrder {
		t.Fatalf(
			`[ERROR] dbapi_test.TestGetAllTasks: Result Length of dbGetAllTasks()=>taskWrapper.Orders = %d, want %d.`,
			actualTasksLengthOrders,
			wantedTasksLengthOrder,
		)
	}

	// Iteration check for parent tasks and subtasks
	i := 0
	// errorMarker := false
	for docID, task := range taskWrapper.Data {
		// // # TEST ordering is correct: [1-5]
		// if task.RankOrder != i+1 {
		// 	t.Fatalf(`[ERROR] dbapi_test.TestGetAllTasks: Ordering for task "%s" = %d, want %d`, task.Content, task.RankOrder, i+1)
		// 	// errorMarker = true
		// }

		// # TEST default attr ID: wanted not "" (empty string)
		if task.ID == "" || task.ID != docID {
			t.Fatalf(`[ERROR] dbapi_test.TestGetAllTasks: Default value of the attribute 'ID' for task "%s" = %s, want %s`, task.Content, task.ID, docID)
		}

		// # TEST default attr IsChecked: wanted false
		if task.IsChecked {
			t.Fatalf(`[ERROR] dbapi_test.TestGetAllTasks: Default value of the attribute 'IsChecked' for task "%s" = %t, want false`, task.Content, task.IsChecked)
		}

		i++
	}

	// Iteration check for parent tasks only
	for i, docID := range taskWrapper.Orders {
		// # TEST ordering is correct: [1-5]
		task := taskWrapper.Data[docID]
		if task.RankOrder != i+1 {
			t.Fatalf(`[ERROR] dbapi_test.TestGetAllTasks-Parent: ordering for task "%s" = %d, want %d`, task.Content, task.RankOrder, i+1)
		}
	}
}

// =================== DANGER ZONE ========================
// make sure target collectionID is correct!
func _deleteCollection(ctx *gin.Context, client *firestore.Client,
	ref *firestore.CollectionRef, batchSize int) error {

	for {
		// Get a batch of documents
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to a
		// WriteBatch.
		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}

// =================== DANGER ZONE ========================
// make sure target collectionID is correct!
func Test__DeleteTestingCollection(t *testing.T) {
	log.Println("[INFO] dbapi_test._DeleteTestingCollection: Deleting all testing documents and collection...")
	testCollectionID := os.Getenv("FIRESTORE_TEST_DATA_COLLECTION_ID")
	ctx := _getTestContext()
	db := _initDbConnection(ctx)
	testCollection := db.Collection(testCollectionID)
	err := _deleteCollection(ctx, db, testCollection, 32)
	if err != nil {
		t.Fatal("[ERROR] DeleteTestingCollection: Failed at deleting test collection")
	}
	db.Close()
}

// TODO: Create test for empty task list

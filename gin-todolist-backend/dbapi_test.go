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

func Test__PreTestingDatabaseCleanup(t *testing.T) {
	err := _CleanUpHelper()
	if err != nil {
		t.Fatalf(`[ERROR] Test__PostTestingDatabaseCleanup: Failed at deleting test collection. Reason: %s`, err)
	}
}

func TestGetAllTasksIsEmpty(t *testing.T) {
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

	wantedTasksLengthOrder := 0
	actualTasksLengthOrders := len(taskWrapper.Orders)
	// - check len(taskWrapper.Orders) == 0
	if actualTasksLengthOrders != wantedTasksLengthOrder {
		t.Fatalf(
			`[ERROR] dbapi_test.TestGetAllTasks: Result Length of dbGetAllTasks()=>taskWrapper.Orders = %d, want %d.`,
			actualTasksLengthOrders,
			wantedTasksLengthOrder,
		)
	}

	wantedTasksLengthData := 0
	actualTasksLengthData := len(taskWrapper.Data)
	// - check len(taskWrapper.Data) == 0
	if actualTasksLengthData != wantedTasksLengthData {
		t.Fatalf(
			`[ERROR] dbapi_test.TestGetAllTasks: Result Length of dbGetAllTasks()=>taskWrapper.Data = %d, want %d.`,
			actualTasksLengthData,
			wantedTasksLengthData,
		)
	}
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

// TODO: Create TestEditTaskContent
func TestEditTaskContent(t *testing.T) {
	c := _getTestContext()
	collectionID := os.Getenv("FIRESTORE_TEST_DATA_COLLECTION_ID")

	taskWrapper, _ := dbGetAllTasks(c, collectionID)
	taskID := taskWrapper.Orders[1]
	oldTask := taskWrapper.Data[taskID]
	newContent := oldTask.Content + "==EDITED VERSION"
	payload := map[string]interface{}{"Content": newContent}

	_, err := dbEditTask(c, taskID, payload, collectionID)
	if err != nil {
		t.Fatalf(`[ERROR] TestEditTaskContent: Failed at uploading edited content. Reason: %s`, err)
	}

	updatedTaskWrapper, _ := dbGetAllTasks(c, collectionID)

	// check if content is updated
	updatedTask := updatedTaskWrapper.Data[taskID]
	if updatedTask.Content != newContent {
		t.Fatalf(
			`[ERROR] TestEditTaskContent: Content is not updated on DB. updatedTask.Content = "%s", want "%s"`,
			updatedTask.Content,
			newContent,
		)
	}

	// check if values of all other fields is not changed or not missing
	if updatedTask.ID != oldTask.ID {
		t.Fatalf(
			`[ERROR] TestEditTaskContent: ID is changed or missing, updatedTask.ID = "%s", want "%s".`,
			updatedTask.ID,
			oldTask.ID,
		)
	}

	if updatedTask.ParentID != oldTask.ParentID {
		t.Fatalf(
			`[ERROR] TestEditTaskContent: ParentID is changed or missing, updatedTask.ID = "%s", want "%s".`,
			updatedTask.ParentID,
			oldTask.ParentID,
		)
	}

	if updatedTask.IsChecked != oldTask.IsChecked {
		t.Fatalf(
			`[ERROR] TestEditTaskContent: IsChecked is changed or missing, updatedTask.IsChecked = "%t", want "%t".`,
			updatedTask.IsChecked,
			oldTask.IsChecked,
		)
	}

	if updatedTask.RankOrder != oldTask.RankOrder {
		t.Fatalf(
			`[ERROR] TestEditTaskContent: RankOrder is changed or missing, updatedTask.RankOrder = "%d", want "%d".`,
			updatedTask.RankOrder,
			oldTask.RankOrder,
		)
	}

	// check whether subtask is equal by length 8
	if len(updatedTask.Subtasks) != len(oldTask.Subtasks) {
		t.Fatalf(
			`[ERROR] TestEditTaskContent: length Subtasks is not equal, len(updatedTask.Subtasks) = "%d", want "%d".`,
			len(updatedTask.Subtasks),
			len(oldTask.Subtasks),
		)
	}

	for i, _ := range updatedTask.Subtasks {
		if updatedTask.Subtasks[i] != oldTask.Subtasks[i] {
			t.Fatalf(
				`[ERROR] TestEditTaskContent: Subtasks is not equal, updatedTask.Subtasks[%d] = "%s", want "%s".`,
				i,
				updatedTask.Subtasks[i],
				oldTask.Subtasks[i],
			)
		}
	}
}

func _TestInvalidEditTaskContent_FieldNotExists(t *testing.T) {
	c := _getTestContext()

	collectionID := os.Getenv("FIRESTORE_TEST_DATA_COLLECTION_ID")
	taskWrapper, _ := dbGetAllTasks(c, collectionID)
	taskID := taskWrapper.Orders[1]
	payload := map[string]interface{}{"InvalidFields": "any value"}

	test, err := dbEditTask(c, taskID, payload, collectionID)
	log.Println(test)
	log.Println(err)
	if err == nil {
		t.Fatalf(
			`[ERROR] TestInvalidEditTaskContent_FieldNotExists: Can still edit task while it should not`,
		)
	}
}

func TestInvalidEditTaskContent_DocumentNotExists(t *testing.T) {
	c := _getTestContext()

	collectionID := os.Getenv("FIRESTORE_TEST_DATA_COLLECTION_ID")
	taskID := "Non-existing taskID (document ID)"
	payload := map[string]interface{}{"InvalidFields": "any value"}

	_, err := dbEditTask(c, taskID, payload, collectionID)
	if err == nil {
		t.Fatalf(
			`[ERROR] TestInvalidEditTaskContent_DocumentNotExists: Can still edit task while document is not exists`,
		)
	}
}

// TODO: Create TestCheckUncheckTask
// TODO: Create TestDeleteSubTask --> delete - auto reoder subtasks (after deleted item)
// TODO: Create TestDeleteParentTask --> delete - auto reoder parents (after deleted item)
// TODO: Create TestAddSubtask

func Test__PostTestingDatabaseCleanup(t *testing.T) {
	err := _CleanUpHelper()
	if err != nil {
		t.Fatalf(`[ERROR] Test__PostTestingDatabaseCleanup: Failed at deleting test collection. Reason: %s`, err)
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
// make sure target collectionID is correct! ==> FIRESTORE_TEST_DATA_COLLECTION_ID
func _CleanUpHelper() error {
	log.Println("[INFO] dbapi_test._CleanUpHelper: Deleting all testing documents and collection...")
	testCollectionID := os.Getenv("FIRESTORE_TEST_DATA_COLLECTION_ID")
	ctx := _getTestContext()
	db := _initDbConnection(ctx)
	testCollection := db.Collection(testCollectionID)
	err := _deleteCollection(ctx, db, testCollection, 128)
	if err != nil {
		db.Close()
		return err
	}
	db.Close()
	return nil
}

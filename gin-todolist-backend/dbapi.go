//dbapi.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func _generateUuid() string {
	newUuid := uuid.New()
	return newUuid.String()
}

func _generateNewTask(content string, parentID string, rankOrder int) Task {
	tmpID := _generateUuid()
	newTask := Task{
		ID:        tmpID,
		ParentID:  parentID,
		RankOrder: rankOrder,
		Content:   content,
		IsChecked: false,
		Subtasks:  nil,
	}
	log.Printf(`[INFO] _generateNewTask: Created task payload--"%s"`, content)
	return newTask
}

func _initFirebaseApp(c *gin.Context) (*firebase.App, error) {
	sa := option.WithCredentialsFile(os.Getenv("RELATIVE_PATH_GCP_KEYFILE"))
	app, err := firebase.NewApp(c, nil, sa)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return app, err
}

func _initFirestoreClient(c *gin.Context, app *firebase.App) (*firestore.Client, error) {
	db, err := app.Firestore(c)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return db, nil
}

func _initDbConnection(c *gin.Context) *firestore.Client {
	app, _ := _initFirebaseApp(c)
	db, _ := _initFirestoreClient(c, app)
	return db
}

func dbCreateTask(c *gin.Context, content string, rankOrder int, collectionID string) (*firestore.WriteResult, error) {
	db := _initDbConnection(c)
	newTask := _generateNewTask(content, "", rankOrder)
	wr, err := db.Collection(collectionID).Doc(newTask.ID).Set(c, newTask)
	if err != nil {
		log.Fatalf(`[Error] dbapi.dbGetAllTasks: Can't Create Task Document with content "%s". %s`, content, err)
		db.Close()
		return nil, err
	}
	db.Close()
	return wr, nil
}

func dbGetAllTasks(c *gin.Context, collectionID string) (*TaskWrapper, error) {
	db := _initDbConnection(c)

	var tasksData = make(map[string]Task)
	var orders []string
	collection := db.Collection(collectionID)
	iter := collection.OrderBy("RankOrder", firestore.Asc).Documents(c)
	counter := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			log.Printf(`[INFO] dbapi.GetAllTasks: Queried %d document(s)`, counter)
			break
		}

		if err != nil {
			log.Fatalln(err)
			db.Close()
			return nil, err
		}

		var task Task
		if err := doc.DataTo(&task); err != nil {
			log.Fatalf(`[Error] dbapi.dbGetAllTasks: Error while converting query data to Task model for doc ID: %s. %v`, doc.Ref.ID, err)
			db.Close()
			return nil, err
		}

		// append tasksData
		tasksData[doc.Ref.ID] = task
		// append orders
		if task.ParentID == "" {
			orders = append(orders, doc.Ref.ID)
		}
		counter++
	}
	db.Close()
	// create taskwrapper
	taskWrapper := TaskWrapper{Data: tasksData, Orders: orders}
	return &taskWrapper, nil
}

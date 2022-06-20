//dbapi.go
package main

import (
  "log"
	"os"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "google.golang.org/api/iterator"

  firebase "firebase.google.com/go"
  "google.golang.org/api/option"
	firestore "cloud.google.com/go/firestore"
)

func _generateUuid() string {
	newUuid := uuid.New()
	return newUuid.String()
}

func _generateNewTask(content string, parentID string, rankOrder int) Task {
	tmpID := _generateUuid()
	newTask := Task{
		ID: tmpID,
		ParentID: parentID,
		RankOrder: rankOrder,
		Content: content,
		IsChecked: false,
		Subtasks: nil,
	}
	return newTask
}

func _initFirebaseApp(c *gin.Context) (*firebase.App, error) {
	sa := option.WithCredentialsFile(os.Getenv("RELATIVE_PATH_GCP_KEYFILE"))
	app, err := firebase.NewApp(c, nil, sa)
	if (err != nil) {
		log.Fatalln(err)
		return nil, err
	}
	return app, err
}

func _initFirestoreClient(c *gin.Context, app *firebase.App) (*firestore.Client, error) {
	db, err := app.Firestore(c)
	if err != nil {
		log.Fatalln(err)
	}
	return db, err
}

func _initDbConnection(c * gin.Context) *firestore.Client {
	app, _ := _initFirebaseApp(c)
	db, _ := _initFirestoreClient(c, app)
	return db
}

func dbCreateTask(c *gin.Context, content string, rankOrder int, collectionID string) (*firestore.WriteResult, error) {
	db := _initDbConnection(c)
	newTask := _generateNewTask(content, "", rankOrder)
	wr, err := db.Collection(collectionID).Doc(newTask.ID).Set(c, newTask)
	if err != nil {
		log.Fatalf(`Database Error: Can't Create Task Document with content "%s"`, content)
		return nil, err
	}
	return wr, nil 
}

// func dbGetAllTasks(context *gin.Context) []Task{
	
// }
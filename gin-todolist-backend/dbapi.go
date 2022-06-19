//dbapi.go
package main

import (
  "log"
	"os"
	// "fmt"
	
	"github.com/gin-gonic/gin"
	// "google.golang.org/api/iterator"

  firebase "firebase.google.com/go"
  "google.golang.org/api/option"
	firestore "cloud.google.com/go/firestore"
)

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
	client, err := app.Firestore(c)
	if err != nil {
		log.Fatalln(err)
	}
	return client, err
}

// func dbGetAllTasks(context *gin.Context) []Task{
	
// }
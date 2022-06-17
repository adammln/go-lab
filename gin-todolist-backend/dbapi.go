//dbapi.go
package main

import (
  "log"
	"os"
	"context"
	"fmt"
	"google.golang.org/api/iterator"

  firebase "firebase.google.com/go"
  "google.golang.org/api/option"
)

// Use a service account

func dbGetAllTasks() {
	ctx := context.Background()
	sa := option.WithCredentialsFile(os.Getenv("RELATIVE_PATH_GCP_KEYFILE"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
	
	log.Print("INITIATING DB CONNECTION...")
	client, err := app.Firestore(ctx)
	log.Print(client)
	if err != nil {
		log.Fatalln(err)
	}

	iter := client.Collection("todolist-app").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
						break
		}
		if err != nil {
						log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}
}
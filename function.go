// Package p contains a Cloud Function that processes Firebase
// Authentication events.
package p

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/functions/metadata"

	firebase "firebase.google.com/go"
	//	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
)

// AuthEvent is the payload of a Firebase Auth event.
// Please refer to the docs for additional information
// regarding Firebase Auth events.
type AuthEvent struct {
	Email string `json:"email"`
	UID   string `json:"uid"`
}

// HelloAuth handles changes to Firebase Auth user objects.
func HelloAuth(ctx context.Context, e AuthEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Function triggered by change to: %v", meta.Resource)
	log.Printf("%v", e)

	// Use the application default credentials
	ctx = context.Background()
	projectID := os.Getenv("GCP_PROJECT")
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	addDocWithID(ctx, client, e)

	return nil

}

func addDocWithID(ctx context.Context, client *firestore.Client, e AuthEvent) error {

	_, err := client.Collection("users").Doc(e.UID).Set(ctx, map[string]interface{}{
		"email" : e.Email,
		"create-timestamp": firestore.ServerTimestamp,
	})
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}

	return err
}

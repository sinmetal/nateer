package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/spanner"
	"github.com/sinmetal/nateer"
	metadatabox "github.com/sinmetalcraft/gcpbox/metadata"
)

func main() {
	log.Print("starting server...")

	ctx := context.Background()

	gcProjectID, err := metadatabox.ProjectID()
	if err != nil {
		panic(err)
	}

	ds, err := datastore.NewClient(ctx, gcProjectID)
	if err != nil {
		panic(err)
	}

	sampleDSStore, err := nateer.NewSampleDSStore(ctx, ds)
	if err != nil {
		panic(err)
	}

	spa, err := spanner.NewClient(ctx, fmt.Sprintf("projects/%s/instances/%s/databases/%s", "gcpug-public-spanner", "merpay-sponsored-instance", "sinmetal"))
	if err != nil {
		panic(err)
	}
	userStore, err := nateer.NewUserStore(ctx, spa)
	if err != nil {
		panic(err)
	}

	externalSendRequestHandler, err := nateer.NewExternalSendRequestHandler(ctx, sampleDSStore, userStore)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/send", externalSendRequestHandler.Handler)
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

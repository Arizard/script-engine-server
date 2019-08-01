package main

import (
	// "context"
	"fmt"
	// "cloud.google.com/go/storage"
	// "io/ioutil"
	"github.com/arizard/script-engine-server/auth"
	"github.com/arizard/script-engine-server/presenters"
	"github.com/arizard/script-engine-server/handlers"
	"github.com/arizard/script-engine-server/infrastructure"
	"net/http"
	"github.com/gorilla/mux"
	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"cloud.google.com/go/firestore"
	"log"
)

func main() {
	fmt.Printf("Setting up infrastructure...\n")

	fmt.Printf("Initialising Firebase Admin SDK...\n")

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
			fmt.Errorf("error initializing app: %v", err)
	}

	projectID := "scriptengine-f031b"

	// Get a Firestore client.
	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Close client when done.
	defer firestoreClient.Close()

	authClient, err := app.Auth(context.Background())
	if err != nil {
			fmt.Printf("error getting Auth client: %v\n", err)
	}

	fmt.Printf("Setting up layers...\n")

	jsonPresenter := presenters.JSONPresenter{}
	userValidator := auth.FirebaseUserValidator{app, authClient}
	docRepo := infrastructure.FirebaseDocumentRepository{
		firestoreClient,
	}

	r := mux.NewRouter().StrictSlash(false)
	JSONHandler := handlers.Handler{
		//DocumentRepo,
		"application/json",
		jsonPresenter,
		userValidator,
		docRepo,
	}

	r.NotFoundHandler = http.HandlerFunc(JSONHandler.NotFoundHandler)

	r.HandleFunc("/", JSONHandler.IndexHandler).Methods("GET")
	r.HandleFunc("/documents", JSONHandler.ListDocumentsHandler).Methods("GET")
	r.HandleFunc("/document/{uuid}", JSONHandler.GetDocumentHandler).Methods("GET")
	r.HandleFunc("/document/", JSONHandler.CreateDocumentHandler).Methods("POST")
	r.HandleFunc("/document/{uuid}", JSONHandler.UpdateDocumentHandler).Methods("PUT")
	r.HandleFunc("/document/{uuid}", JSONHandler.DeleteDocumentHandler).Methods("DELETE")

	// r.HandleFunc("/look/{name}", JSONHandler.GetPublicURL).Methods("GET")
	// r.HandleFunc("/give", JSONHandler.UploadUserFile).Methods("POST")

	// fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	// r.PathPrefix("/static/").Handler(fs)

	http.ListenAndServe(":8080", r)


}
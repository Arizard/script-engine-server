package main

import (
	// "context"
	"fmt"
	// "cloud.google.com/go/storage"
	// "io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/arizard/script-engine-server/auth"
	"github.com/arizard/script-engine-server/handlers"
	"github.com/arizard/script-engine-server/infrastructure"
	"github.com/arizard/script-engine-server/presenters"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func main() {
	fmt.Printf("Setting up infrastructure...\n")

	fmt.Printf("Initialising Firebase Admin SDK...\n")

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}

	projectID := "probable-spoon"

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

	fmt.Println("Check")

	queryDefaultDocument, err := firestoreClient.Collection("documents").Doc("default-document").Get(context.Background())

	if err != nil {
		log.Fatalf("An error occured. %v\n", err)
	}
	mapping := queryDefaultDocument.Data()

	defaultDocumentData := mapping["data"].(string)

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
		defaultDocumentData,
	}

	r.NotFoundHandler = http.HandlerFunc(JSONHandler.NotFoundHandler)

	r.HandleFunc("/", JSONHandler.IndexHandler).Methods("GET")
	r.HandleFunc("/documents", JSONHandler.CORSWrapper(JSONHandler.ListDocumentsHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/document/{uuid}", JSONHandler.CORSWrapper(JSONHandler.GetDocumentHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/document/", JSONHandler.CORSWrapper(JSONHandler.CreateDocumentHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/document/{uuid}", JSONHandler.CORSWrapper(JSONHandler.UpdateDocumentHandler)).Methods("PUT", "OPTIONS")
	r.HandleFunc("/document/{uuid}", JSONHandler.DeleteDocumentHandler).Methods("DELETE", "OPTIONS")

	// r.HandleFunc("/look/{name}", JSONHandler.GetPublicURL).Methods("GET")
	// r.HandleFunc("/give", JSONHandler.UploadUserFile).Methods("POST")

	// fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	// r.PathPrefix("/static/").Handler(fs)

	http.ListenAndServe(":8080", r)

}

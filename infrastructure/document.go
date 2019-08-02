package infrastructure

import (
	"google.golang.org/api/iterator"
	// "fmt"
	"log"
	"golang.org/x/net/context"
	"github.com/arizard/script-engine-server/entities"
	"github.com/google/uuid"
	"cloud.google.com/go/firestore"
)

type FirebaseDocumentRepository struct {
	Client *firestore.Client
}


func (repo FirebaseDocumentRepository) List(uid string) []entities.Document {
	docs := []entities.Document{}

	docsQuery := repo.Client.Collection("scriptengine-documents").Where("owner", "==", uid).OrderBy("title", firestore.Asc).Documents(context.Background())

	for {
		doc, err := docsQuery.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error while iterating query: %s", err)
			break
		}
		mapping := doc.Data()
		docs = append(
			docs,
			entities.Document{
				Name: mapping["name"].(string),
				Title: mapping["title"].(string),
				Owner: mapping["owner"].(string),
				Data: mapping["data"].(string),
				UUID: mapping["uuid"].(string),
			},
		)
	}

	return docs
}


func (repo FirebaseDocumentRepository) Add(doc entities.Document) string {
	
	_uuid, _ := uuid.NewRandom()

	doc.UUID = _uuid.String()

	_, err := repo.Client.Collection("scriptengine-documents").Doc(_uuid.String()).Set(
		context.Background(),
		map[string]interface{}{
			"name": doc.Name,
			"title": doc.Title,
			"data": doc.Data,
			"owner": doc.Owner,
			"uuid": _uuid.String(),
		},
	)
	if err != nil {
		log.Printf("error while writing to cloud firestore: %s", err)
	}
	return _uuid.String()
}

func (repo FirebaseDocumentRepository) Get(_uuid string) (*entities.Document, error) {
	query, err := repo.Client.Collection("scriptengine-documents").Doc(_uuid).Get(context.Background())
	if err != nil {
		log.Printf("error while reading from cloud firestore: %s", err)
		return nil, err
	}

	mapping := query.Data()

	return &entities.Document{
		Name: mapping["name"].(string),
		Title: mapping["title"].(string),
		Owner: mapping["owner"].(string),
		Data: mapping["data"].(string),
		UUID: mapping["uuid"].(string),
	}, nil
}

func (repo FirebaseDocumentRepository) Update(_uuid string, doc *entities.Document) {
	_, err := repo.Client.Collection("scriptengine-documents").Doc(_uuid).Set(
		context.Background(),
		map[string]interface{}{
			"name": doc.Name,
			"title": doc.Title,
			"data": doc.Data,
			"owner": doc.Owner,
			"uuid": _uuid,
		},
	)
	if err != nil {
		log.Printf("error while writing to cloud firestore: %s", err)
	}
	return
}

func (repo FirebaseDocumentRepository) Delete(_uuid string) {
	return
}



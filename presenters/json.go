package presenters

import (
	"github.com/arizard/script-engine-server/entities"
	// "bytes"
	// "text/template"
	"fmt"
	"encoding/json"
)

// All the JSON presenter should do is take a struct as an argument, and then
// output the correct JSON formatted string result.

// JSONPresenter presents data in browser-renderable html.
type JSONPresenter struct {}

// NotFound presents the 404 output.
func (JSONPresenter) NotFound() string {
	return "{ message: \"Not Found\" }"
}


// Forbidden presents the 403 output.
func (JSONPresenter) Forbidden() string {
	return "{ message: \"Forbidden\" }"
}

// InternalServerError displays the 500 output.
func (JSONPresenter) InternalServerError() string {
	return "{ message: \"Internal Server Error.\" }"
}

// Index displays the index output.
func (JSONPresenter) Index() string {
	return "{ message: \"Index\" }"
}

// ListDocuments displays the document list output.
func (JSONPresenter) ListDocuments(docs []entities.Document) string {
	mapping := [](map[string]interface{}){}

	for _, doc := range docs {
		mapping = append(
			mapping,
			map[string]interface{}{
				"name": doc.Name,
				"title": doc.Title,
				"owner": doc.Owner,
				"data": "(omitted by presenter)",
				"uuid": doc.UUID,
			},
		)
	}

	data, err := json.Marshal(mapping)

	if err != nil {
		return fmt.Sprintf("An error occured while marshalling: %s", err)
	}

	return string(data)
}

// CreateDocument creates and saves a new document, returning the uuid.
func (JSONPresenter) CreateDocument(uuid string) string {
	mapping := map[string]interface{}{
		"uuid": uuid,
	}

	data, err := json.Marshal(mapping)

	if err != nil {
		return fmt.Sprintf("An error occured while marshalling: %s", err)
	}

	return string(data)
}

// GetDocument retrieves the record of a document by the uuid.
func (JSONPresenter) GetDocument(doc entities.Document) string {
	mapping := map[string]interface{}{
		"name": doc.Name,
		"title": doc.Title,
		"owner": doc.Owner,
		"data": doc.Data,
		"uuid": doc.UUID,
	}

	data, err := json.Marshal(mapping)

	if err != nil {
		return fmt.Sprintf("An error occured while marshalling: %s", err)
	}

	return string(data)
}

func (JSONPresenter) UpdateDocument() string {
	return "{ message: \"UpdateDocument\" }"
}

func (JSONPresenter) DeleteDocument() string {
	return "{ message: \"DeleteDocument\" }"
}


package presenters

import (
	"github.com/arizard/script-engine-server/entities"
)

// Presenter defines the contract for presenters, either html or json.
type Presenter interface {
	NotFound() string
	InternalServerError() string
	Forbidden() string
	Index() string
	ListDocuments([]entities.Document) string
	CreateDocument(string) string
	GetDocument(entities.Document) string
	UpdateDocument() string
	DeleteDocument() string
}

package entities

import (
	
)

// Document is the entity which describes a user-generated file.
type Document struct {
	Name string // Name must be unique.
	Title string
	Owner string
	Data string // JSON formatted
	UUID string
}

// DocumentService manages business logic related to 
// instances of Document.
type DocumentService struct {
	Repository DocumentRepository
}

// NewDocument creates a new Document instance.
func (s *DocumentService) NewDocument(name string, title string, data string) Document {
	return Document{
		Name: name,
		Title: title,
		Data: data,
	}
}

// DocumentRepository manages the storage of Document instances.
type DocumentRepository interface {
	List(userUUID string) []Document
	Add(Document Document) string
	Get(_uuid string) (*Document, error)
	Update(_uuid string, Document *Document)
	Delete(_uuid string)
}

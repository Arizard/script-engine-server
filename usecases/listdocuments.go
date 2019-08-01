package usecases

import (
	// "io/ioutil"
	"github.com/arizard/script-engine-server/entities"
	// "mime/multipart"
)

type ListDocuments struct {
	DocumentRepository entities.DocumentRepository
	DocumentOwner string
	PresenterFunc func ([]entities.Document) string
	Response *ResponseCollector
}

func (uc ListDocuments) Setup() {

}

func (uc ListDocuments) Execute() {

	docs := uc.DocumentRepository.List(uc.DocumentOwner)
	jsonOut := uc.PresenterFunc(docs)

	resp := Response{
		Body: jsonOut,
	}

	uc.Response.SetResponse(&resp)

}
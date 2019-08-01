package usecases

import (
	// "io/ioutil"
	"github.com/arizard/script-engine-server/entities"
	// "mime/multipart"
)

type CreateDocument struct {
	DocumentRepository entities.DocumentRepository
	DocumentName string
	DocumentTitle string
	DocumentOwner string
	DocumentData string
	Response *ResponseCollector
}

func (uc CreateDocument) Setup() {

}

func (uc CreateDocument) Execute() {

	newDocument := entities.Document{
		uc.DocumentName,
		uc.DocumentTitle,
		uc.DocumentOwner,
		uc.DocumentData,
		"",
	}

	uuid := uc.DocumentRepository.Add(newDocument)

	resp := Response{
		Body: uuid,
	}

	uc.Response.SetResponse(&resp)

}
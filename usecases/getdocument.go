package usecases

import (
	"fmt"
	// "io/ioutil"
	"github.com/arizard/script-engine-server/entities"
	// "mime/multipart"
)

type GetDocument struct {
	DocumentRepository entities.DocumentRepository
	DocumentOwner string
	DocumentUUID string
	PresenterFunc func (entities.Document) string
	Response *ResponseCollector
}

func (uc GetDocument) Setup() {

}

func (uc GetDocument) Execute() {

	doc, err := uc.DocumentRepository.Get(uc.DocumentUUID)

	if err != nil {
		resp := ResponseError{
			"REPOSITORY_FAILURE",
			"Failed to call Get on the repository.",
		}
		uc.Response.SetError(&resp)
		fmt.Printf("Failed get %s", err)
		return
	}

	if (doc.Owner != uc.DocumentOwner) {
		resp := ResponseError{
			"FORBIDDEN",
			"User is not allowed to access this document.",
		}
		uc.Response.SetError(&resp)
		fmt.Printf("Failed get %s", err)
		return
	}

	jsonOut := uc.PresenterFunc(*doc)

	resp := Response{
		Body: jsonOut,
	}

	uc.Response.SetResponse(&resp)

}

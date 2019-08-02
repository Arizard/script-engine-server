package usecases

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"github.com/arizard/script-engine-server/entities"
	// "mime/multipart"
)

type UpdateDocument struct {
	DocumentRepository entities.DocumentRepository
	DocumentOwner string
	DocumentUUID string
	DocumentJSON []byte
	PresenterFunc func () string
	Response *ResponseCollector
}

func (uc UpdateDocument) Setup() {

}

func (uc UpdateDocument) Execute() {

	mapping := map[string]interface{}{}
	json.Unmarshal(uc.DocumentJSON, &mapping)

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

	JSONData, _ := json.Marshal(mapping["data"])
	doc.Data = string(JSONData)
	doc.Title = mapping["title"].(string)

	uc.DocumentRepository.Update(uc.DocumentUUID, doc)

	//jsonOut := uc.PresenterFunc(*doc)

	resp := Response{
		Body: "{ message: \"Success\"}",
	}

	uc.Response.SetResponse(&resp)

}

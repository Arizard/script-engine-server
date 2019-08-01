package usecases

import (
	"io/ioutil"
	"github.com/arizard/fish-less-coffee/entities"
	"mime/multipart"
)

type CreateUserFile struct {
	File multipart.File
	FileHeader *multipart.FileHeader
	UserFileRepo entities.UserFileRepository
	Response *ResponseCollector
}

func (uc CreateUserFile) Setup() {

}

func (uc CreateUserFile) Execute() {
	defer panicHandler(uc.Response)
	userFileService := entities.UserFileService{
		Repository: uc.UserFileRepo,
	}
	data, err := ioutil.ReadAll(uc.File)
	if err != nil {
		respErr := ResponseError{
			Name: "GENERIC_FAILURE",
			Description: "Something went wrong.",
		}

		uc.Response.SetError(&respErr)
		return
	}
	newUserFile := userFileService.NewUserFile(uc.FileHeader.Filename, data)

	uc.UserFileRepo.Add(&newUserFile)

	resp := Response{
		Body: newUserFile.Name,
	}

	uc.Response.SetResponse(&resp)

}
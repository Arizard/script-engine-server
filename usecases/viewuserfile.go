package usecases

import (
	"github.com/arizard/fish-less-coffee/presenters"
	"github.com/arizard/fish-less-coffee/entities"
)

type ViewUserFile struct {
	FileName string
	UserFileRepo entities.UserFileRepository
	Presenter presenters.Presenter
	Response *ResponseCollector
}

func (uc ViewUserFile) Setup() {

}

func (uc ViewUserFile) Execute() {
	defer panicHandler(uc.Response)
	resp := Response{
		Body: uc.Presenter.GetUserFile(
			uc.FileName,
			uc.UserFileRepo.GetPublicURL(uc.FileName),
		),
	}

	uc.Response.SetResponse(&resp)
}
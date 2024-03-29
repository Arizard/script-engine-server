package handlers

import (
	"io/ioutil"
	"log"
	// "encoding/json"
	"github.com/arizard/script-engine-server/presenters"
	"github.com/arizard/script-engine-server/auth"
	"github.com/arizard/script-engine-server/entities"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"github.com/arizard/script-engine-server/usecases"
	"github.com/harlow/authtoken"
)

// The handler should handle HTTP headers and status codes, and execute the
// use cases, then run the output through the presenter.

// Handler is a struct which implements methods that take the 
// ResponseWriter and Request objects as arguments, such as from an
// HTTP request. It is used to decouple the Drivers layer from the
// Controllers and Presenters.
type Handler struct {
	//DocumentRepo entities.DocumentRepository
	ContentType string
	Presenter presenters.Presenter
	UserValidator auth.UserValidator
	DocumentRepository entities.DocumentRepository
	DefaultDocumentData string
}

func (handler Handler) VerifyRequest(r *http.Request, verifyFunc func(string) bool) bool {
	token, err := authtoken.FromRequest(r)
	if (err != nil){
		token = ""
	}
	return verifyFunc(token) 
}

func (handler Handler) CORSWrapper(hf func (http.ResponseWriter, *http.Request)) func (http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", handler.ContentType)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control")
		w.Header().Set("Access-Control-Expose-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if (r.Method == "OPTIONS") {
			log.Printf("incoming CORS request")
		} else {
			log.Printf("incoming request")
			w.Header().Set("Content-Type", handler.ContentType)
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization")
			hf(w, r)
		}
	}
	
}


// NotFoundHandler handles 404s
func (handler Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", handler.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(404)
	fmt.Fprintf(w, handler.Presenter.NotFound())
}

// ForbiddenHandler handles 403s
func (handler Handler) ForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", handler.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(403)
	fmt.Fprintf(w, handler.Presenter.Forbidden())
}

// InternalServerErrorHandler handles 500s
func (handler Handler) InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", handler.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(500)
	fmt.Fprintf(w, handler.Presenter.InternalServerError())
}

// IndexHandler handles a request for the Index view of the presenter.
func (handler Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", handler.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, handler.Presenter.Index())
}

// ListDocumentsHandler handles requests to the list document endpoint
func (handler Handler) ListDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	verified := handler.VerifyRequest(r, handler.UserValidator.CanListDocuments)
	if (verified == false){
		handler.ForbiddenHandler(w, r)
		log.Printf("user is forbidden to access.")
		return
	}

	idToken, _ := authtoken.FromRequest(r)
	uid := handler.UserValidator.GetUserRef(idToken)

	rc := usecases.ResponseCollector{}
	uc := usecases.ListDocuments{
		handler.DocumentRepository,
		uid,
		handler.Presenter.ListDocuments,	
		&rc,
	}

	uc.Setup()
	uc.Execute()

	fmt.Fprintf(w, rc.Response.Body)
}

// CreateDocumentHandler handles request to the create document endpoint.
func (handler Handler) CreateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	verified := handler.VerifyRequest(r, handler.UserValidator.CanCreateDocument)
	if (verified == false){
		handler.ForbiddenHandler(w, r)
		return
	}

	idToken, _ := authtoken.FromRequest(r)
	uid := handler.UserValidator.GetUserRef(idToken)

	rc := usecases.ResponseCollector{}
	uc := usecases.CreateDocument{
		handler.DocumentRepository,
		"no_name",
		"New Document",
		uid,
		handler.DefaultDocumentData,	
		&rc,
	}

	uc.Setup()
	uc.Execute()
	
	w.Header().Set("Content-Type", handler.ContentType)
	fmt.Fprintf(w, handler.Presenter.CreateDocument(rc.Response.Body))
}

func (handler Handler) GetDocumentHandler(w http.ResponseWriter, r *http.Request) {
	verified := handler.VerifyRequest(r, handler.UserValidator.CanGetDocument)
	if (verified == false){
		handler.ForbiddenHandler(w, r)
		return
	}

	idToken, _ := authtoken.FromRequest(r)
	uid := handler.UserValidator.GetUserRef(idToken)

	docUUID := mux.Vars(r)["uuid"]

	rc := usecases.ResponseCollector{}
	uc := usecases.GetDocument{
		handler.DocumentRepository,
		uid,
		docUUID,
		handler.Presenter.GetDocument,
		&rc,
	}

	uc.Setup()
	uc.Execute()

	if rc.Error != nil {
		handler.NotFoundHandler(w, r)
		return
	}
	
	fmt.Fprintf(w, rc.Response.Body)
}

func (handler Handler) UpdateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	verified := handler.VerifyRequest(r, handler.UserValidator.CanUpdateDocument)
	if (verified == false){
		handler.ForbiddenHandler(w, r)
		return
	}

	idToken, _ := authtoken.FromRequest(r)
	uid := handler.UserValidator.GetUserRef(idToken)

	docUUID := mux.Vars(r)["uuid"]

	docJSON, _ := ioutil.ReadAll(r.Body)

	rc := usecases.ResponseCollector{}
	uc := usecases.UpdateDocument{
		handler.DocumentRepository,
		uid,
		docUUID,
		docJSON,
		handler.Presenter.UpdateDocument,
		&rc,
	}

	uc.Setup()
	uc.Execute()

	if rc.Error != nil {
		handler.InternalServerErrorHandler(w, r)
		return
	}
	
	fmt.Fprintf(w, rc.Response.Body)
}

func (handler Handler) DeleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	verified := handler.VerifyRequest(r, handler.UserValidator.CanDeleteDocument)
	if (verified == false){
		handler.ForbiddenHandler(w, r)
		return
	}
	w.Header().Set("Content-Type", handler.ContentType)
	fmt.Fprintf(w, handler.Presenter.DeleteDocument())
}

// // GetPublicURL handles the GetPublicURL view of the presenter.
// func (handler Handler) GetPublicURL(w http.ResponseWriter, r *http.Request) {
// 	fileName := mux.Vars(r)["name"]

// 	rc := usecases.ResponseCollector{}
// 	uc := usecases.ViewUserFile{
// 		FileName: fileName,
// 		UserFileRepo: handler.UserFileRepo,
// 		Presenter: handler.Presenter,
// 		Response: &rc,
// 	}

// 	uc.Setup()
// 	uc.Execute()
	
// 	if rc.Error != nil {
// 		if rc.Error.Name == "NOT_FOUND" {
// 			w.WriteHeader(404)
// 			handler.NotFound(w, r)
// 		}
// 		if rc.Error.Name == "SEVERE_FAILURE" {
// 			w.WriteHeader(500)
// 			handler.InternalServerErrorHandler(w, r)
// 		}
// 		return
// 	}
// 	fmt.Fprintf(w, "%s", rc.Response.Body)
// }


// func (handler Handler) UploadUserFile(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseMultipartForm(32 << 20)
// 	if err != nil {
// 		fmt.Printf("%s\n", err)
// 	}
// 	file, fileHeader, _ := r.FormFile("file")

// 	rc := usecases.ResponseCollector{}
// 	uc := usecases.CreateUserFile{
// 		File: file,
// 		FileHeader: fileHeader,
// 		DocumentRepo: handler.DocumentRepo,
// 		Response: &rc,
// 	}

// 	uc.Setup()

// 	uc.Execute()

// 	w.Header().Set("Content-Location", fmt.Sprintf("/look/%s", rc.Response.Body))
// 	http.Redirect(w, r, fmt.Sprintf("/look/%s", rc.Response.Body), 301)

// 	handler.errorHelper(w, r, rc)
// }
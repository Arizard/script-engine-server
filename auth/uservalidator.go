package auth

import (
	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"firebase.google.com/go/auth"
	"fmt"
)

// UserValidator is an interface defining the contract for user validators.
type UserValidator interface {
	CanListDocuments(token string) bool
	CanCreateDocument(token string) bool
	CanGetDocument(token string) bool
	CanUpdateDocument(token string) bool
	CanDeleteDocument(token string) bool
	GetUserRef(token string) string
}

// FirebaseUserValidator is a UserValidator which uses Firebase to validate user's ID token.
type FirebaseUserValidator struct {
	App *firebase.App
	Client *auth.Client
}

func (fuv FirebaseUserValidator) GetUserRef(idToken string) string {
	result, err := fuv.Client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
			fmt.Printf("error verifying ID token: %v\n", err)
			return "nil"
	}

	return result.UID
}

func (fuv FirebaseUserValidator) VerifyToken(idToken string) bool {
	// fmt.Printf("Bearer: %s\n", idToken)

	result, err := fuv.Client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
			fmt.Printf("error verifying ID token: %v\n", err)
			return false
	}

	fmt.Printf("Verified UID %v.\n", result.UID)

	return true
}

func (fuv FirebaseUserValidator) CanListDocuments(token string) bool {
	return fuv.VerifyToken(token)
}

func (fuv FirebaseUserValidator) CanDirectoryList(token string) bool {
	return fuv.VerifyToken(token)
}

func (fuv FirebaseUserValidator) CanCreateDocument(token string) bool {
	return fuv.VerifyToken(token)
}

func (fuv FirebaseUserValidator) CanGetDocument(token string) bool {
	return fuv.VerifyToken(token)
}

func (fuv FirebaseUserValidator) CanUpdateDocument(token string) bool {
	return fuv.VerifyToken(token)
}

func (fuv FirebaseUserValidator) CanDeleteDocument(token string) bool {
	return fuv.VerifyToken(token)
}
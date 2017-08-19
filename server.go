package main

import (
	"log"
	// Standard library packages
	"net/http"

	"github.com/44r0n/SessionManager/helpers"
	"github.com/44r0n/SessionManager/repository"

	"github.com/44r0n/SessionManager/controllers"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Instantiate a new router
	r := httprouter.New()
	const serverURL = "127.0.0.1:3000"
	connString := helpers.GetConnString("configuration/configuration.json")

	// Get a UserController instance
	repo, err := repository.NewUserRepository(connString)
	if err != nil {
		log.Fatalf("Cannot load user repository: %v", err)
	}
	uc := controllers.NewUserController(repo)
	r.POST("/Register", uc.Register)
	r.POST("/Login", uc.Login)
	r.POST("/Logout", uc.Logout)
	r.POST("/Token/isValid", uc.CheckToken)

	log.Printf("Starting server at %v", serverURL)
	// Fire up the server
	log.Fatal(http.ListenAndServe(serverURL, r))
}

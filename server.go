package main

import (
	"log"
	// Standard library packages
	"net/http"

	// Third party packages
	"SessionManager/controllers"
	"SessionManager/helpers"
	"SessionManager/models/user"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Instantiate a new router
	r := httprouter.New()
	const serverURL = "localhost:3000"
	connString := helpers.GetConnString("configuration/configuration.json")

	// Get a UserController instance
	uc := controllers.NewUserController(models.NewUserRepository(connString))
	r.POST("/Register", uc.Register)
	r.POST("/Login", uc.Login)
	r.POST("/Logout", uc.Logout)
	r.POST("/Token/isValid", uc.CheckToken)

	log.Printf("Starting server at %v", serverURL)
	// Fire up the server
	log.Fatal(http.ListenAndServe(serverURL, r))
}

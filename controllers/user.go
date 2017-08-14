package controllers

import (
	"SessionManager/helpers"
	"SessionManager/models/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// UserController represents the controller for operating on the User resource
type UserController struct {
	userRepo models.IUserRepositoryInterface
}

// NewUserController creates UserController
func NewUserController(UserRepo models.IUserRepositoryInterface) *UserController {
	usc := new(UserController)
	usc.userRepo = UserRepo
	return usc
}

// Register function to register an user recieved in json format
func (uc *UserController) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Printf("/Register")
	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		uc.handleError(w, err)
		return
	}

	existsUsername, e := uc.userRepo.ExistsUsername(u.UserName)
	if e != nil {
		uc.handleError(w, e)
		return
	}

	existsEmail, e2 := uc.userRepo.ExistsEmail(u.Email)
	if e2 != nil {
		uc.handleError(w, e2)
		return
	}

	if existsUsername || existsEmail {
		var maperrors map[string]string
		maperrors = make(map[string]string)
		if existsUsername {
			maperrors["username"] = "An account already exists with this username"
		}
		if existsEmail {
			maperrors["email"] = "An account already exists with this email"
		}

		httpError := helpers.HTTPErrorHandler{
			Status:      http.StatusConflict,
			Error:       "FIELDS_REPEATED",
			Description: "One or more fields already exist",
			Fields:      maperrors,
		}
		uc.responseError(w, httpError)
		return
	}
	if err := uc.userRepo.Register(u); err != nil {
		uc.handleError(w, err)
		return

	}
	w.WriteHeader(http.StatusCreated)
}

//Login controller function
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Printf("/Login")
	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		uc.handleError(w, err)
		return
	}
	token, err := uc.userRepo.LogIn(u.UserName, u.Password)
	if err != nil {
		uc.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if token != "" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"response":{"status":"OK","token":"`+token+`","error":""}}`)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, `{"response":{"status":"Incorrect user or password","token":"","error":""}}`)
}

func (uc *UserController) responseError(w http.ResponseWriter, httpError helpers.HTTPErrorHandler) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpError.Status)
	httpe, e := json.Marshal(httpError)
	if e != nil {
		uc.handleError(w, e)
		return
	}
	fmt.Fprintf(w, `{"error":%s}`, httpe)
}

func (uc *UserController) handleError(w http.ResponseWriter, e error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{"error": {"status":500,"error":"UNCONTROLLED_ERROR", "description":`+e.Error()+`}}`)
}

//Logout controller function
func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Printf("/Logout")
	token := uc.checkTokenHeader(w, r)
	if token == "" {
		return
	}

	if err := uc.userRepo.LogOut(token); err != nil {
		httpError := helpers.HTTPErrorHandler{
			Status:      http.StatusNotFound,
			Error:       "FIELDS_REPEATED",
			Description: "Invalid token",
			Fields:      nil,
		}
		uc.responseError(w, httpError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (uc *UserController) checkTokenHeader(w http.ResponseWriter, r *http.Request) string {
	arraytoken, exists := r.Header["Authorization"]
	if !exists {
		httpError := helpers.HTTPErrorHandler{
			Status:      http.StatusInternalServerError,
			Error:       "INVALID_TOKEN",
			Description: "There must be a valid Authorization token",
			Fields:      nil,
		}
		uc.responseError(w, httpError)
		return ""
	}

	if len(arraytoken) == 0 {
		httpError := helpers.HTTPErrorHandler{
			Status:      http.StatusInternalServerError,
			Error:       "INVALID_TOKEN",
			Description: "There must be a valid Authorization token",
			Fields:      nil,
		}
		uc.responseError(w, httpError)
		return ""
	}

	if arraytoken[0] == "" {
		httpError := helpers.HTTPErrorHandler{
			Status:      http.StatusNotFound,
			Error:       "INVALID_TOKEN",
			Description: "There must be a valid Authorization token",
			Fields:      nil,
		}
		uc.responseError(w, httpError)
		return ""
	}

	return arraytoken[0]
}

//CheckToken controller function
func (uc *UserController) CheckToken(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Printf("/Token/isValid")
	token := uc.checkTokenHeader(w, r)
	result, err := uc.userRepo.CheckToken(token)
	if err != nil {
		httpError := helpers.HTTPErrorHandler{
			Status:      http.StatusInternalServerError,
			Error:       "Error",
			Description: "Unexpected error",
			Fields:      nil,
		}
		uc.responseError(w, httpError)
		return
	}

	if result {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

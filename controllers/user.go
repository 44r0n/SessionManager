package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/44r0n/SessionManager/helpers"

	"github.com/44r0n/SessionManager/codes"
	"github.com/44r0n/SessionManager/models"
	"github.com/44r0n/SessionManager/repository"

	"github.com/julienschmidt/httprouter"
)

// UserController represents the controller for operating on the User resource
type UserController struct {
	userRepo repository.IUserRepositoryInterface
}

// NewUserController creates UserController
func NewUserController(UserRepo repository.IUserRepositoryInterface) UserController {
	if UserRepo == nil {
		log.Fatal("UserRepo cannot be nil")
	}
	usc := new(UserController)
	usc.userRepo = UserRepo
	return *usc
}

// Register function to register an user recieved in json format
func (uc *UserController) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	response := models.Response{Error: codes.Unknown}
	responseData := models.ResponseData{Data: response}
	log.Printf("/Register")
	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response = models.Response{Status: http.StatusBadRequest,
			Error:       codes.JSonError,
			Description: "Failed decoding json"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		log.Printf("Failed decoding json: %v", err)
		return
	}

	if u.UserName == "" || u.Password == "" || u.Email == "" {
		response = models.Response{Status: http.StatusBadRequest,
			Error:       codes.JSonError,
			Description: "Some params required are empty"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		log.Printf("Some required params are empty: %v", err)
		return
	}

	existsUsername, e := uc.userRepo.ExistsUsername(u.UserName)
	if e != nil {
		response = models.Response{Status: http.StatusInternalServerError,
			Error:       codes.DataBaseError,
			Description: "There was an error with the database"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		log.Printf("Falied cheching if exists username. Error: %v", e)
		return
	}

	if existsUsername {
		response = models.Response{Status: http.StatusConflict,
			Error:       codes.RepeatedUserName,
			Description: "Repeated UserName"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		return
	}

	existsEmail, e2 := uc.userRepo.ExistsEmail(u.Email)
	if e2 != nil {
		response = models.Response{Status: http.StatusInternalServerError,
			Error:       codes.DataBaseError,
			Description: "There was an error with the database"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		log.Printf("Falied cheching if exists email. Error: %v", e)
		return
	}

	if existsEmail {
		response = models.Response{Status: http.StatusConflict,
			Error:       codes.RepeatedUserEmail,
			Description: "Repeated Email"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		return
	}

	if err := uc.userRepo.Register(u); err != nil {
		response = models.Response{Status: http.StatusInternalServerError,
			Error:       codes.DataBaseError,
			Description: "There was an error with the database"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		log.Printf("Error registering user: %v. Error: %v", u.UserName, err)
		return

	}

	response = models.Response{Status: http.StatusCreated, Error: codes.Ok}
	responseData.Data = response
	uc.responseToClient(w, responseData)

}

//Login controller function
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	const passError = "crypto/bcrypt: hashedPassword is not the hash of the given password"
	response := models.Response{Error: codes.Unknown}
	responseData := models.ResponseData{Data: response}
	log.Printf("/Login")
	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response = models.Response{Status: http.StatusInternalServerError,
			Error:       codes.JSonError,
			Description: "Failed decoding json"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		log.Printf("Failed decoding json: %v", err)
		return
	}

	userID, pass, err := uc.userRepo.GetIDAndPassword(u.UserName)
	if err != nil {
		response = models.Response{Status: http.StatusInternalServerError,
			Error:       codes.DataBaseError,
			Description: "There was an error with the database"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		log.Printf("Failed loging in: %v", err)
		return
	}

	if pass == "" {
		response = models.Response{Status: http.StatusNotFound,
			Error:       codes.UserNotFound,
			Description: "User not found"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		return
	}

	err = helpers.CheckHash(pass, u.Password)
	if err != nil {
		response = models.Response{Status: http.StatusNotFound,
			Error:       codes.UserNotFound,
			Description: "User not found"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		log.Printf("Failed checking password in: %v", err)
		return
	}

	token, err := helpers.Tokenize(userID)
	if err != nil {
		log.Printf("Failed generating token: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if token != "" {
		err = uc.userRepo.CreateToken(userID, token)
		if err != nil {
			response = models.Response{Status: http.StatusInternalServerError,
				Error:       codes.DataBaseError,
				Description: "There was an error with the database"}
			responseData.Data = response
			uc.responseToClient(w, responseData)
			log.Printf("Failed creating token: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"Response":{"Status":`+strconv.Itoa(http.StatusOK)+`,"Token":"`+token+`","Error":`+strconv.Itoa(codes.Ok)+`}}`)
		return
	}

	response = models.Response{Status: http.StatusNotFound,
		Error:       codes.UserNotFound,
		Description: "User not found"}
	responseData.Data = response
	uc.responseToClient(w, responseData)
}

func (uc *UserController) responseToClient(w http.ResponseWriter, response models.ResponseData) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Data.Status)
	json, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed decodign response to client. Error: %v", err)
		return
	}
	fmt.Fprintf(w, string(json[:]))
}

//Logout controller function
func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	response := models.Response{Error: codes.Unknown}
	responseData := models.ResponseData{Data: response}
	log.Printf("/Logout")
	token := uc.checkTokenHeader(w, r)
	if token == "" {
		response = models.Response{Status: http.StatusBadRequest,
			Error:       codes.NoTokenProvided,
			Description: "No token was provided"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		return
	}

	userID, err := helpers.GetFromToken(token)
	if err != nil {
		response = models.Response{Status: http.StatusInternalServerError,
			Error:       codes.InvalidToken,
			Description: "The token is invalid"}
		return
	}

	if err := uc.userRepo.DeleteToken(userID, token); err != nil {
		response = models.Response{Status: http.StatusInternalServerError,
			Error:       codes.DataBaseError,
			Description: "There was an error with the database"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		log.Printf("Error loging out %v", err)
		return
	}

	response = models.Response{Status: http.StatusOK,
		Error:       codes.Ok,
		Description: ""}
	responseData.Data = response
	uc.responseToClient(w, responseData)
}

func (uc *UserController) checkTokenHeader(w http.ResponseWriter, r *http.Request) string {
	arraytoken, exists := r.Header["Authorization"]
	if !exists {
		return ""
	}

	if len(arraytoken) == 0 {
		return ""
	}

	return arraytoken[0]
}

//CheckToken controller function
func (uc *UserController) CheckToken(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	response := models.Response{Error: codes.Unknown}
	responseData := models.ResponseData{Data: response}
	log.Printf("/Token/isValid")
	token := uc.checkTokenHeader(w, r)
	if token == "" {
		response = models.Response{Status: http.StatusBadRequest,
			Error:       codes.NoTokenProvided,
			Description: "No token was provided"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		return
	}
	result, err := uc.userRepo.CheckToken(token)
	if err != nil {
		response = models.Response{Status: http.StatusInternalServerError,
			Error:       codes.DataBaseError,
			Description: "There was an error with the database"}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		log.Printf("Error Checking Token: %v", err)
		return
	}

	if result {
		response = models.Response{Status: http.StatusOK,
			Error:       codes.Ok,
			Description: ""}
		responseData.Data = response
		uc.responseToClient(w, responseData)
		return
	}

	response = models.Response{Status: http.StatusNotFound,
		Error:       codes.InvalidToken,
		Description: "The token is invalid"}
	responseData.Data = response
	uc.responseToClient(w, responseData)
}

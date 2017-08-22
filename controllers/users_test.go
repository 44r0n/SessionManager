package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/44r0n/SessionManager/codes"
	"github.com/44r0n/SessionManager/helpers"
	"github.com/44r0n/SessionManager/models"
	"github.com/44r0n/SessionManager/repository"

	"github.com/julienschmidt/httprouter"
	. "github.com/smartystreets/goconvey/convey"
)

type UserRepositoryTest struct {
	err        error
	validUser  bool
	validEmail bool
	token      string
}

func (usrt *UserRepositoryTest) Register(user models.User) error {
	return usrt.err
}

func (usrt *UserRepositoryTest) LogIn(userName, password string) (string, error) {
	return usrt.token, usrt.err
}

func (usrt *UserRepositoryTest) LogOut(token string) error {
	return usrt.err
}

func (usrt *UserRepositoryTest) ExistsUsername(userName string) (bool, error) {
	return usrt.validUser, usrt.err
}

func (usrt *UserRepositoryTest) ExistsEmail(email string) (bool, error) {
	return usrt.validEmail, usrt.err
}

func (usrt *UserRepositoryTest) CheckToken(token string) (bool, error) {
	return usrt.validUser, usrt.err
}

func NewUserRepositoryTest(user, email bool, errs error, token string) repository.IUserRepositoryInterface {
	usrt := UserRepositoryTest{errs, user, email, token}
	return &usrt
}

var (
	database   = flag.Bool("database", false, "run database integration tests")
	connString string
)

func TestMain(m *testing.M) {
	flag.Parse()
	if *database {
		log.Printf("Running tests with database integration")
		setupDatabase()
	} else {
		log.Printf("Running tests with NO database integration")
	}

	retCode := m.Run()
	teardownFunction()
	os.Exit(retCode)
}

func setupDatabase() {
	configFile := "../configuration/configuration.json"

	connString = helpers.GetConnString(configFile)
	if connString == "" {
		log.Fatalf("Connection string is empty")
	}
	log.Printf("The database has been set up wit connection string: %v", connString)
}

func teardownFunction() {

}

func TestRegisterHandler(t *testing.T) {
	Convey("Given a valid username", t, func() {
		var jsonStr = []byte(`{"UserName":"Bob Smith","Email":"mail@mail.com","Password":"secretPassword"}`)

		Convey("It should be registered. With status 201", func() {
			var uc UserController
			if *database {
				repo, err := repository.NewUserRepository(connString)
				if err != nil {
					t.Fatal(err)
				}
				uc = NewUserController(repo)
			} else {
				uc = NewUserController(NewUserRepositoryTest(false, false, nil, ""))
			}
			req, err := http.NewRequest("POST", "/Register", bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := httprouter.New()

			router.Handle("POST", "/Register", uc.Register)
			router.ServeHTTP(rr, req)
			status := rr.Code
			So(status, ShouldEqual, http.StatusCreated)
			response := models.ResponseData{}
			err = json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Failed unmarshaling response: %v", err)
			}
			So(response.Data.Status, ShouldEqual, http.StatusCreated)
			So(response.Data.Error, ShouldEqual, codes.Ok)
			So(response.Data.Description, ShouldBeEmpty)
		})
	})
}

func TestRegisterBadJon(t *testing.T) {
	Convey("Given an invalid json to register, it should return a bad request", t, func() {
		var jsonStr = []byte(`{"user":{"username":"Aaron","email":"correo@mail.com","password":"secretP1assword"}}`)
		uc := NewUserController(NewUserRepositoryTest(false, false, nil, ""))
		req, err := http.NewRequest("POST", "/Register", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := httprouter.New()

		router.Handle("POST", "/Register", uc.Register)
		router.ServeHTTP(rr, req)
		status := rr.Code
		So(status, ShouldEqual, http.StatusBadRequest)
		response := models.ResponseData{}
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed unmarshaling response: %v", err)
		}
		So(response.Data.Status, ShouldEqual, http.StatusBadRequest)
		So(response.Data.Error, ShouldEqual, codes.JSonError)
		So(response.Data.Description, ShouldEqual, "Some params required are empty")
	})
}

func TestRegisterRepeatedUser(t *testing.T) {
	Convey("Given a registered user", t, func() {
		var repo repository.IUserRepositoryInterface
		var err error
		if *database {
			repo, err = repository.NewUserRepository(connString)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			repo = NewUserRepositoryTest(false, false, nil, "")
		}

		rr := registerUser([]byte(`{"UserName":"Repeated","Email":"repeat@mail.com","Password":"secretPassword"}`),
			repo, t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusCreated)

		Convey("Cannot register it with the same username and email", func() {
			rr := registerUser([]byte(`{"UserName":"Repeated","Email":"repeat@mail.com","Password":"secretPassword"}`),
				NewUserRepositoryTest(true, true, nil, ""), t)

			status := rr.Code
			So(status, ShouldEqual, http.StatusConflict)

			response := models.ResponseData{}
			err = json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Failed unmarshaling response: %v", err)
			}
			So(response.Data.Status, ShouldEqual, http.StatusConflict)
			So(response.Data.Error, ShouldEqual, codes.RepeatedUserName)
			So(response.Data.Description, ShouldEqual, "Repeated UserName")
		})
	})
}

func registerUser(user []byte, userRepo repository.IUserRepositoryInterface, t *testing.T) *httptest.ResponseRecorder {
	uc := NewUserController(userRepo)

	req, err := http.NewRequest("POST", "/Register", bytes.NewBuffer(user))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := httprouter.New()

	router.Handle("POST", "/Register", uc.Register)
	router.ServeHTTP(rr, req)

	return rr
}

func TestRegisterRepeatedUserName(t *testing.T) {
	Convey("Given a registered user", t, func() {
		var repo repository.IUserRepositoryInterface
		var err error
		if *database {
			repo, err = repository.NewUserRepository(connString)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			repo = NewUserRepositoryTest(false, false, nil, "")
		}
		rr := registerUser([]byte(`{"UserName":"RepeatedUsername","Email":"repeatusernam@mail.com","Password":"secretPassword"}`),
			repo, t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusCreated)

		Convey("Cannot register user with the same username", func() {
			var repo repository.IUserRepositoryInterface
			var err error
			if *database {
				repo, err = repository.NewUserRepository(connString)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				repo = NewUserRepositoryTest(true, false, nil, "")
			}
			rr := registerUser([]byte(`{"UserName":"RepeatedUsername","Email":"other@mail.com","Password":"secretPassword"}`),
				repo, t)

			status := rr.Code
			So(status, ShouldEqual, http.StatusConflict)

			response := models.ResponseData{}
			err = json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Failed unmarshaling response: %v", err)
			}
			So(response.Data.Status, ShouldEqual, http.StatusConflict)
			So(response.Data.Error, ShouldEqual, codes.RepeatedUserName)
			So(response.Data.Description, ShouldEqual, "Repeated UserName")
		})
	})
}

func TestRegisterRepeatedMail(t *testing.T) {
	Convey("Given a registered user", t, func() {
		var repo repository.IUserRepositoryInterface
		var err error
		if *database {
			repo, err = repository.NewUserRepository(connString)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			repo = NewUserRepositoryTest(false, false, nil, "")
		}
		rr := registerUser([]byte(`{"UserName":"RepeatedUsername2","Email":"repeatusernam2@mail.com","Password":"secretPassword"}`),
			repo, t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusCreated)

		Convey("Cannot register user with the same email", func() {
			var repo repository.IUserRepositoryInterface
			var err error
			if *database {
				repo, err = repository.NewUserRepository(connString)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				repo = NewUserRepositoryTest(false, true, nil, "")
			}
			rr := registerUser([]byte(`{"UserName":"RepeatedMail2","Email":"repeatusernam2@mail.com","Password":"secretPassword"}`),
				repo, t)

			status := rr.Code
			So(status, ShouldEqual, http.StatusConflict)

			response := models.ResponseData{}
			err = json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Failed unmarshaling response: %v", err)
			}
			So(response.Data.Status, ShouldEqual, http.StatusConflict)
			So(response.Data.Error, ShouldEqual, codes.RepeatedUserEmail)
			So(response.Data.Description, ShouldEqual, "Repeated Email")
		})
	})
}

func TestRegisterUnexpectedError(t *testing.T) {
	Convey("Given a valid user and no database connection", t, func() {
		rr := registerUser([]byte(`{"UserName":"RepeatedUsername","Email":"repeatusernam@mail.com","Password":"secretPassword"}`),
			NewUserRepositoryTest(false, false, errors.New("No bd connection"), ""), t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusInternalServerError)

		response := models.ResponseData{}
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed unmarshaling response: %v", err)
		}
		So(response.Data.Status, ShouldEqual, http.StatusInternalServerError)
		So(response.Data.Error, ShouldEqual, codes.DataBaseError)
		So(response.Data.Description, ShouldEqual, "There was an error with the database")
	})
}

func TestLoginOK(t *testing.T) {
	Convey("Given a valid user, it can log in", t, func() {
		var repo repository.IUserRepositoryInterface
		var err error
		if *database {
			repo, err = repository.NewUserRepository(connString)
			if err != nil {
				panic(err)
			}

			user := models.User{
				UserName: "LogOK",
				Email:    "logok@mail.com",
				Password: "secretPassword",
			}
			err = repo.Register(user)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			repo = NewUserRepositoryTest(true, false, nil, "1234abcd")
		}

		rr := simulateLogin(&repo, []byte(`{"UserName":"LogOK","Email":"","Password":"secretPassword"}`), t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusOK)
		if !*database {
			expected := `{"Response":{"Status":` + strconv.Itoa(http.StatusOK) + `,"Token":"1234abcd","Error":` + strconv.Itoa(codes.Ok) + `}}`
			So(rr.Body.String(), ShouldEqual, expected)
		}
	})
}

func TestLoginBadPassword(t *testing.T) {
	Convey("Given an invalid Password of a registered user, it should return not found", t, func() {
		var repo repository.IUserRepositoryInterface
		var err error
		if *database {
			repo, err = repository.NewUserRepository(connString)
			if err != nil {
				panic(err)
			}

			user := models.User{
				UserName: "LogOK2",
				Email:    "logok2@mail.com",
				Password: "secretPassword",
			}
			err = repo.Register(user)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			repo = NewUserRepositoryTest(false, false, nil, "")
		}

		rr := simulateLogin(&repo, []byte(`{"UserName":"LogOK2","Email":"","Password":"secretPassworda"}`), t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusNotFound)

		response := models.ResponseData{}
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed unmarshaling response: %v", err)
		}
		So(response.Data.Status, ShouldEqual, http.StatusNotFound)
		So(response.Data.Error, ShouldEqual, codes.UserNotFound)
		So(response.Data.Description, ShouldEqual, "User not found")
	})
}

func simulateLogin(usrt *repository.IUserRepositoryInterface, jsonUser []byte, t *testing.T) *httptest.ResponseRecorder {
	uc := NewUserController(*usrt)
	req, err := http.NewRequest("POST", "/Login", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := httprouter.New()

	router.Handle("POST", "/Login", uc.Login)
	router.ServeHTTP(rr, req)
	return rr
}

func TestLoginNotOK(t *testing.T) {
	Convey("Given a invalid user, it cannot log in", t, func() {
		var repo repository.IUserRepositoryInterface
		var err error
		if *database {
			repo, err = repository.NewUserRepository(connString)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			repo = NewUserRepositoryTest(false, false, nil, "")
		}
		rr := simulateLogin(&repo, []byte(`{"UserName":"Bob Smitha","Email":"","Password":"secretPassword"}`), t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusNotFound)

		response := models.ResponseData{}
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed unmarshaling response: %v", err)
		}
		So(response.Data.Status, ShouldEqual, http.StatusNotFound)
		So(response.Data.Error, ShouldEqual, codes.UserNotFound)
		So(response.Data.Description, ShouldEqual, "User not found")
	})
}

func TestLoginErrorDB(t *testing.T) {
	Convey("Given a valid user, and fails to connect", t, func() {
		repo := NewUserRepositoryTest(false, false, errors.New("No bd connection"), "")
		rr := simulateLogin(&repo, []byte(`{"UserName":"Bob Smitha","Password":"secretPassword"}`), t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusInternalServerError)

		response := models.ResponseData{}
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed unmarshaling response: %v", err)
		}
		So(response.Data.Status, ShouldEqual, http.StatusInternalServerError)
		So(response.Data.Error, ShouldEqual, codes.DataBaseError)
		So(response.Data.Description, ShouldEqual, "There was an error with the database")
	})
}

func TestLogoutOK(t *testing.T) {
	Convey("Given a valid token, it should log out", t, func() {
		var repo repository.IUserRepositoryInterface
		var token string
		var err error
		if *database {
			repo, err = repository.NewUserRepository(connString)
			if err != nil {
				t.Fatal(err)
			}
			user := models.User{
				UserName: "TestLogout",
				Email:    "logout@mail.com",
				Password: "passlogout",
			}
			repo.Register(user)
			token, err = repo.LogIn(user.UserName, user.Password)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			repo = NewUserRepositoryTest(true, false, nil, "")
			token = "1234abcd"
		}
		rr := simulateLogout(repo, token, t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusOK)

		response := models.ResponseData{}
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed unmarshaling response: %v", err)
		}
		So(response.Data.Status, ShouldEqual, http.StatusOK)
		So(response.Data.Error, ShouldEqual, codes.Ok)
		So(response.Data.Description, ShouldBeEmpty)
	})
}

func simulateLogout(usrt repository.IUserRepositoryInterface, token string, t *testing.T) *httptest.ResponseRecorder {
	uc := NewUserController(usrt)
	req, err := http.NewRequest("POST", "/Logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := httprouter.New()

	router.Handle("POST", "/Logout", uc.Logout)
	router.ServeHTTP(rr, req)
	return rr
}

func TestLogoutNotOK(t *testing.T) {
	Convey("Given an invalid token, it should do nothing", t, func() {
		rr := simulateLogout(NewUserRepositoryTest(false, false, nil, ""), "", t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusBadRequest)

		response := models.ResponseData{}
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed unmarshaling response: %v", err)
		}
		So(response.Data.Status, ShouldEqual, http.StatusBadRequest)
		So(response.Data.Error, ShouldEqual, codes.NoTokenProvided)
		So(response.Data.Description, ShouldEqual, "There was no token provided")
	})
}

func TestCheckTokenOK(t *testing.T) {
	Convey("Given a valid token, it shold return ok when it is checked", t, func() {
		repo := NewUserRepositoryTest(true, false, nil, "1234abcd")
		rr := simulateCheckToken(&repo, "1234abcd", t)
		status := rr.Code
		So(status, ShouldEqual, http.StatusOK)

		response := models.ResponseData{}
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed unmarshaling response: %v", err)
		}
		So(response.Data.Status, ShouldEqual, http.StatusOK)
		So(response.Data.Error, ShouldEqual, codes.Ok)
		So(response.Data.Description, ShouldBeEmpty)
	})
}

func simulateCheckToken(usrt *repository.IUserRepositoryInterface, token string, t *testing.T) *httptest.ResponseRecorder {
	uc := NewUserController(*usrt)
	req, err := http.NewRequest("POST", "/Token/isValid", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := httprouter.New()

	router.Handle("POST", "/Token/isValid", uc.CheckToken)
	router.ServeHTTP(rr, req)
	return rr
}

func TestCeckTokenNotOK(t *testing.T) {
	Convey("Given an invalid token, it should return not found when it is checked", t, func() {
		repo := NewUserRepositoryTest(false, false, nil, "")
		rr := simulateCheckToken(&repo, "45jvm", t)
		status := rr.Code
		So(status, ShouldEqual, http.StatusNotFound)

		response := models.ResponseData{}
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed unmarshaling response: %v", err)
		}
		So(response.Data.Status, ShouldEqual, http.StatusNotFound)
		So(response.Data.Error, ShouldEqual, codes.InvalidToken)
		So(response.Data.Description, ShouldEqual, "The token is invalid")
	})
}

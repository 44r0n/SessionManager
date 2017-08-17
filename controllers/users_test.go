package controllers

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/44r0n/SessionManager/helpers"
	models "github.com/44r0n/SessionManager/models/user"

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

func NewUserRepositoryTest(user, email bool, errs error, token string) *UserRepositoryTest {
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
	ip := helpers.GetIP(configFile)
	cmdStr := "../initsql.sh ''" + ip
	cmd := exec.Command("/bin/sh", cmdStr)
	_, err := cmd.Output()

	if err != nil {
		log.Fatalf(err.Error())
		return
	}

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
		var jsonStr = []byte(`{"name":"Bob Smith","email":"mail@mail.com","pass":"secretPassword"}`)

		Convey("It should be registered. With status 201", func() {
			var uc *UserController
			if *database {
				repo, err := models.NewUserRepository(connString)
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
		})
	})
}

func TestRegisterRepeatedUser(t *testing.T) {
	Convey("Given a registered user", t, func() {
		var repo models.IUserRepositoryInterface
		var err error
		if *database {
			repo, err = models.NewUserRepository(connString)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			repo = NewUserRepositoryTest(false, false, nil, "")
		}

		rr := registerUser([]byte(`{"name":"Repeated","email":"repeat@mail.com","pass":"secretPassword"}`),
			repo, t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusCreated)

		Convey("Cannot register it with the same username and email", func() {
			rr := registerUser([]byte(`{"name":"Repeated","email":"repeat@mail.com","pass":"secretPassword"}`),
				NewUserRepositoryTest(true, true, nil, ""), t)

			status := rr.Code
			So(status, ShouldEqual, http.StatusConflict)

			expected := `{"error":{"status":409,"error":"FIELDS_REPEATED","description":"One or more fields already exist","fields":{"email":"An account already exists with this email","username":"An account already exists with this username"}}}`
			So(expected, ShouldEqual, rr.Body.String())
		})
	})
}

func registerUser(user []byte, userRepo models.IUserRepositoryInterface, t *testing.T) *httptest.ResponseRecorder {
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
		var repo models.IUserRepositoryInterface
		var err error
		if *database {
			repo, err = models.NewUserRepository(connString)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			repo = NewUserRepositoryTest(false, false, nil, "")
		}
		rr := registerUser([]byte(`{"name":"RepeatedUsername","email":"repeatusernam@mail.com","pass":"secretPassword"}`),
			repo, t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusCreated)

		Convey("Cannot register user with the same username", func() {
			var repo models.IUserRepositoryInterface
			var err error
			if *database {
				repo, err = models.NewUserRepository(connString)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				repo = NewUserRepositoryTest(true, false, nil, "")
			}
			rr := registerUser([]byte(`{"name":"RepeatedUsername","email":"other@mail.com","pass":"secretPassword"}`),
				repo, t)

			status := rr.Code
			So(status, ShouldEqual, http.StatusConflict)

			expected := `{"error":{"status":409,"error":"FIELDS_REPEATED","description":"One or more fields already exist","fields":{"username":"An account already exists with this username"}}}`
			So(expected, ShouldEqual, rr.Body.String())
		})
	})
}

func TestRegisterRepeatedMail(t *testing.T) {
	Convey("Given a registered user", t, func() {
		var repo models.IUserRepositoryInterface
		var err error
		if *database {
			repo, err = models.NewUserRepository(connString)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			repo = NewUserRepositoryTest(false, false, nil, "")
		}
		rr := registerUser([]byte(`{"name":"RepeatedUsername2","email":"repeatusernam2@mail.com","pass":"secretPassword"}`),
			repo, t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusCreated)

		Convey("Cannot register user with the same email", func() {
			var repo models.IUserRepositoryInterface
			var err error
			if *database {
				repo, err = models.NewUserRepository(connString)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				repo = NewUserRepositoryTest(false, true, nil, "")
			}
			rr := registerUser([]byte(`{"name":"RepeatedMail2","email":"repeatusernam2@mail.com","pass":"secretPassword"}`),
				repo, t)

			status := rr.Code
			So(status, ShouldEqual, http.StatusConflict)

			expected := `{"error":{"status":409,"error":"FIELDS_REPEATED","description":"One or more fields already exist","fields":{"email":"An account already exists with this email"}}}`
			So(rr.Body.String(), ShouldEqual, expected)
		})
	})
}

func TestRegisterUnexpectedError(t *testing.T) {
	Convey("Given a valid user and no database connection", t, func() {
		rr := registerUser([]byte(`{"name":"RepeatedUsername","email":"repeatusernam@mail.com","pass":"secretPassword"}`),
			NewUserRepositoryTest(false, false, errors.New("No bd connection"), ""), t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestLoginOK(t *testing.T) {
	Convey("Given a valid user, it can log in", t, func() {
		var repo models.IUserRepositoryInterface
		if *database {
			repo, err := models.NewUserRepository(connString)
			if err != nil {
				t.Fatal(err)
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
		rr := simulateLogin(repo, []byte(`{"name":"LogOK","email":"","pass":"secretPassword"}`), t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusOK)
		if !*database {
			expected := `{"response":{"status":"OK","token":"1234abcd","error":""}}`
			So(rr.Body.String(), ShouldEqual, expected)
		}
	})
}

func simulateLogin(usrt models.IUserRepositoryInterface, jsonUser []byte, t *testing.T) *httptest.ResponseRecorder {
	uc := NewUserController(usrt)
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
		var repo models.IUserRepositoryInterface
		var err error
		if *database {
			repo, err = models.NewUserRepository(connString)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			repo = NewUserRepositoryTest(false, false, nil, "")
		}
		rr := simulateLogin(repo, []byte(`{"name":"Bob Smitha","email":"","pass":"secretPassword"}`), t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusNotFound)

		expected := `{"response":{"status":"Incorrect user or password","token":"","error":""}}`
		So(rr.Body.String(), ShouldEqual, expected)
	})
}

func TestLoginErrorDB(t *testing.T) {
	Convey("Given a valid user, and fails to connect", t, func() {
		rr := simulateLogin(NewUserRepositoryTest(false, false, errors.New("No bd connection"), ""), []byte(""), t)

		status := rr.Code
		So(status, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestLogoutOK(t *testing.T) {
	Convey("Given a valid token, it should log out", t, func() {
		var repo models.IUserRepositoryInterface
		var token string
		var err error
		if *database {
			repo, err = models.NewUserRepository(connString)
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
		So(status, ShouldEqual, http.StatusNoContent)
	})
}

func simulateLogout(usrt models.IUserRepositoryInterface, token string, t *testing.T) *httptest.ResponseRecorder {
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
		So(status, ShouldEqual, http.StatusNotFound)
	})
}

func TestCheckTokenOK(t *testing.T) {
	Convey("Given a valid token, it shold return ok when it is checked", t, func() {
		rr := simulateCheckToken(NewUserRepositoryTest(true, false, nil, "1234abcd"), "1234abcd", t)
		status := rr.Code
		So(status, ShouldEqual, http.StatusNoContent)
	})
}

func simulateCheckToken(usrt *UserRepositoryTest, token string, t *testing.T) *httptest.ResponseRecorder {
	uc := NewUserController(usrt)
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
		rr := simulateCheckToken(NewUserRepositoryTest(false, false, nil, ""), "45jvm", t)
		status := rr.Code
		So(status, ShouldEqual, http.StatusNotFound)
	})

}

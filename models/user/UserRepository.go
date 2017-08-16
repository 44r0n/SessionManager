package models

import (
	"fmt"

	"github.com/44r0n/SessionManager/data"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository struct implementation of IUserRepositoryInterface
type UserRepository struct {
	mysqlconnString string
}

// NewUserRepository function to get new UserRepository
func NewUserRepository(connString string) (*UserRepository, error) {
	if connString == "" {
		return nil, fmt.Errorf("connString cannot be void string")
	}
	usr := UserRepository{connString}
	return &usr, nil
}

// Register function that registers the given user.
func (usr *UserRepository) Register(user User) error {
	password := []byte(user.Password)
	hashedPass, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	datab := database.NewDatabaseConnection(usr.mysqlconnString)
	if err = datab.ExecuteNonQuery("INSERT INTO users (id,username,email,password,date_created) VALUES (uuid(),?,?,?,NOW())", user.UserName, user.Email, hashedPass); err != nil {
		return err
	}
	return nil
}

// ExistsUsername function checks if the given userName exists
func (usr *UserRepository) ExistsUsername(userName string) (bool, error) {
	datab := database.NewDatabaseConnection(usr.mysqlconnString)
	rows, err := datab.ExecuteQuery("SELECT  username from users where username = ? LIMIT 1", userName)
	if err != nil {
		return false, err
	}
	var userNameChecker string
	rows.Next()
	rows.Scan(&userNameChecker)
	if userNameChecker != "" {
		return true, nil
	}
	return false, nil
}

// ExistsEmail function checks if the given email exists
func (usr *UserRepository) ExistsEmail(email string) (bool, error) {
	datab := database.NewDatabaseConnection(usr.mysqlconnString)
	rows, err := datab.ExecuteQuery("SELECT  email from users where email = ? LIMIT 1", email)
	if err != nil {
		return false, err
	}
	var emailChecker string
	rows.Next()
	rows.Scan(&emailChecker)
	if emailChecker != "" {
		return true, nil
	}
	return false, nil
}

// LogIn Logs In the user
func (usr *UserRepository) LogIn(userName string, password string) (string, error) {
	datab := database.NewDatabaseConnection(usr.mysqlconnString)
	rows, err := datab.ExecuteQuery("SELECT id, password from users where username = ? LIMIT 1",
		userName)
	if err != nil {
		return "", err
	}
	var idChecker, storedPassword string
	rows.Next()
	rows.Scan(&idChecker, &storedPassword)
	if idChecker == "" {
		return "", nil
	}

	if storedPassword == "" {
		return "", fmt.Errorf("Could not get password from database")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return "", err
	}

	tokenString, err := tokenize(idChecker)
	if err != nil {
		return "", err
	}

	if err = datab.ExecuteNonQuery("INSERT INTO user_tokens (user, token, last_date_used) VALUES (?,?,NOW())", idChecker, tokenString); err != nil {
		return "", err
	}

	return tokenString, nil
}

func tokenize(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id": id,
	})
	tokenString, err := token.SignedString([]byte("SecretKey"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// LogOut function
func (usr *UserRepository) LogOut(token string) error {
	user, err := getUserFromToken(token)
	if err != nil {
		return err
	}
	datab := database.NewDatabaseConnection(usr.mysqlconnString)
	if err := datab.ExecuteNonQuery("DELETE FROM user_tokens where user = ? and token = ?", user, token); err != nil {
		return err
	}
	return nil
}

func getUserFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("SecretKey"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["id"].(string), nil
	}

	return "", err
}

// CheckToken function
func (usr *UserRepository) CheckToken(token string) (bool, error) {
	user, err := getUserFromToken(token)
	if err != nil {
		return false, err
	}

	datab := database.NewDatabaseConnection(usr.mysqlconnString)
	rows, err := datab.ExecuteQuery("SELECT user from user_tokens where user = ? AND token = ?",
		user, token)
	if err != nil {
		return false, err
	}
	var idChecker string
	rows.Next()
	rows.Scan(&idChecker)
	if idChecker == "" {
		return false, nil
	}

	return true, nil
}

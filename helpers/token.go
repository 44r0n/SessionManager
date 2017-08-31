package helpers

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// Tokenize returns a token from a given text
func Tokenize(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id": id,
	})
	tokenString, err := token.SignedString([]byte("SecretKey"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetFromToken gets the value of a given token
func GetFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("SecretKey"), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["id"].(string), nil
	}

	return "", err
}

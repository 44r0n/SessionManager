package helpers

import "golang.org/x/crypto/bcrypt"

//GenerateHash from a given string
func GenerateHash(text string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed[:]), nil
}

//CheckHash from a given hasedText and a Text
func CheckHash(hashedText, text string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(text))
}

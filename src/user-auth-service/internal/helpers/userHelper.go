package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashed)
}

func VerifyingOfPassword(userPassword, foundUserPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(foundUserPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		check = false
		msg = "Email or password is incorrect"
	}
	return check, msg
}

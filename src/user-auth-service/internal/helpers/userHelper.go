package helpers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CheckType(c *gin.Context, userId string) string {
	uid := c.GetString("uid")
	userType := c.GetString("user_type")
	if userType == "ADMIN" && uid != userId {
		return ""
	} else {
		return "Unauthorized access to the server"
	}
}

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

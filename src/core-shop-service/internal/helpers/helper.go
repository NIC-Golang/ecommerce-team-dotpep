package helpers

import "github.com/gin-gonic/gin"

func CheckUserType(c *gin.Context, role string) bool {
	return c.GetString("client_type") == role
}

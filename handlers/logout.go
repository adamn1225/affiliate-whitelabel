package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	// Clear cookie by setting MaxAge to -1
	c.SetCookie("token", "", -1, "/", "", false, true) // Clear cookie by setting MaxAge to -1
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

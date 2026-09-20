package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"sosu-backend/internal/service"
)

func AuthRequired(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak ditemukan, silakan login"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "sesi tidak valid, silakan login ulang"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
package middleware

import (
	"log"
	"net/http"
	"strings"

	 
	"consistent_1/Infrastructure/auth"   

	"github.com/gin-gonic/gin" 
	"consistent_1/Domain"    
)
func AuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, err := jwtService.GetUserIDFromToken(tokenString)
		if err != nil {
			log.Printf("JWT validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidToken.Error()})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
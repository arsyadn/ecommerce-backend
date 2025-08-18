package middleware

import (
	"final-project/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// cek jika header Authorization ada
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		userId, role, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", userId)
		c.Set("user_role", role)
	}
}


func RoleAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleInterface, exists := c.Get("user_role")

		if !exists || roleInterface != "admin" {
			c.JSON(403, gin.H{"error": "Only admin can access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RoleUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists || userRole != "user" {
			c.JSON(403, gin.H{"error": "Only user can access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

const (
	RoleVendor    = "vendor"
	RoleAffiliate = "affiliate"
	RoleAdmin     = "admin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
		return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["role"] != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		c.Next()
	}
}

func RequireRole(expectedRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        var tokenString string

        // Check Authorization header or fall back to cookie
        if strings.HasPrefix(authHeader, "Bearer ") {
            tokenString = strings.TrimPrefix(authHeader, "Bearer ")
        } else if cookieToken, err := c.Cookie("token"); err == nil {
            tokenString = cookieToken
        }

        if tokenString == "" {
            log.Println("Missing or malformed token")
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
            return
        }

        // Log the token for debugging
        log.Printf("Received token: %s", tokenString)

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            log.Printf("Invalid token: %v", err)
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            log.Println("Invalid token structure")
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token structure"})
            return
        }

        roleClaim := claims["role"]
        if roleClaim != expectedRole {
            log.Printf("Access denied for role: %v", roleClaim)
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
            return
        }

        c.Set("user_id", claims["user_id"])
        c.Set("role", roleClaim)
        c.Next()
    }
}


func GenerateToken(userID any, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		// Check Authorization header or fall back to cookie
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else if cookieToken, err := c.Cookie("token"); err == nil {
			tokenString = cookieToken
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token structure"})
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}


// func RequireAnyRole(allowedRoles ...string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := extractToken(c) // your logic here
// 		claims := parseClaims(tokenString)

// 		userRole := claims["role"].(string)
// 		for _, role := range allowedRoles {
// 			if userRole == role {
// 				c.Set("user_id", claims["user_id"])
// 				c.Set("role", userRole)
// 				c.Next()
// 				return
// 			}
// 		}

// 		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
// 	}
// }




package middleware

import (
	"e-commerce/internal/config"
	"e-commerce/internal/utils/xgin"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			xgin.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authorization header required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" || tokenString == authHeader {
			xgin.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid authorization header format")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			xgin.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid or expired token")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			xgin.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid token claims")
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			xgin.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid token claims")
			c.Abort()
			return
		}
		if exp, ok := claims["exp"].(float64); ok {
			expiration := time.Unix(int64(exp), 0)
			if time.Now().After(expiration) {
				xgin.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Token has expired")
				c.Abort()
				return
			}
		}
		c.Set("user_id", userID)
		c.Next()
	}

}

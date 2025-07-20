package main

import (
	"log"
	"strings"
	"test-vk/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/registration" || c.Request.URL.Path == "/login" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tokenString == "" {
			if c.Request.URL.Path == "/list" {
				c.Next()
				return
			} else {
				c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
				return
			}
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatal("Failed to load config:", err)
			}

			return []byte(cfg.Secret.Secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("user_login", claims["user_login"])
		c.Next()
	}
}

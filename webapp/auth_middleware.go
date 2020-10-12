package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var (
	// TODO : generate cryptographic secret and store in the env
	secret = "dev-secret"
)

func CreateToken(id string) string {
	// Create the token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	token.Claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return ""
	}

	return tokenString
}

func ParseToken(authHeader string) (string, error) {
	// split ["Bearer", <token>"]
	bearerToken := strings.Split(authHeader, " ")
	// get the <token>
	tokenVal := bearerToken[len(bearerToken)-1]

	parsedToken, err := jwt.Parse(tokenVal, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if id, ok := claims["id"].(string); ok {
			return id, nil
		} else {
			return "", nil
		}
	} else {
		return "", err
	}
}

func AuthMW(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			b := []byte(secret)
			return b, nil
		})

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
	}
}

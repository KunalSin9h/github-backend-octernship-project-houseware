package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func unAuthorizedResponse(c *gin.Context, err error) {
	sendResponse("UnAuthorized", err.Error(), nil, c, http.StatusUnauthorized)
	c.Abort()
}

func (app *Config) AuthorizationMiddleware(c *gin.Context) {

	authTokenString, err := c.Cookie("Authorization")

	if err != nil {
		unAuthorizedResponse(c, err)
		return
	}

	token, err := jwt.Parse(authTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		unAuthorizedResponse(c, err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["userId"]
		expirationTime := claims["exp"]

		if userId == nil || expirationTime == nil {
			unAuthorizedResponse(c, errors.New("invalid auth token"))
			return
		}

		if float64(time.Now().Unix()) > expirationTime.(float64) {
			unAuthorizedResponse(c, errors.New("auth token expired"))
			return
		}

		c.Set("userId", userId)

		c.Next()
	} else {
		unAuthorizedResponse(c, errors.New("invalid auth token"))
		return
	}
}

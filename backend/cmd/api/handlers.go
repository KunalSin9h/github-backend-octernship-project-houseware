package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type responsePayload struct {
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

type loginRequestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *Config) login(c *gin.Context) {

	dataBytes, err := io.ReadAll(c.Request.Body)

	if err != nil {
		res := responsePayload{
			Message: "Error reading request body",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	var reqPayload loginRequestPayload
	json.Unmarshal(dataBytes, &reqPayload)

	username := reqPayload.Username
	password := reqPayload.Password

	if username == "" || password == "" {
		res := responsePayload{
			Message: "Missing Username or Password in request",
			Error:   "missing username or password in request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, res)
	}

}

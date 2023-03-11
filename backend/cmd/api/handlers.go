package main

import (
	"houseware---backend-engineering-octernship-KunalSin9h/data"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type responsePayload struct {
	Message string         `json:"message"`
	Error   string         `json:"error"`
	Data    map[string]any `json:"data,omitempty"`
}

func (app *Config) login(c *gin.Context) {

	var reqPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := c.Bind(&reqPayload)

	if err != nil {
		res := responsePayload{
			Message: "Error reading request body",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	username := reqPayload.Username
	password := reqPayload.Password

	if username == "" || password == "" {
		res := responsePayload{
			Message: "Missing Username or Password in request",
			Error:   "missing username or password in request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	user, err := app.Models.User.GetByUsername(username)

	if err != nil {
		res := responsePayload{
			Message: "Error while getting user",
			Error:   "error while getting user",
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if user.ID == "" {
		// User Does not exist
		res := responsePayload{
			Message: "Invalid username or password",
			Error:   "invalid username or password",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	isPasswordMatched, err := user.PasswordMatch(password)

	if err != nil {
		res := responsePayload{
			Message: "Error while verifying password",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if !isPasswordMatched {
		// invalid password
		res := responsePayload{
			Message: "Invalid username or password",
			Error:   "invalid username or password",
			Data:    nil,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(JWT_SECRET))

	if err != nil {
		res := responsePayload{
			Message: "Failed to create JWT Token",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	res := responsePayload{
		Message: "User Signed in",
		Error:   "",
		Data: map[string]any{
			"user": user,
		},
	}

	c.JSON(http.StatusOK, res)
}

func (app *Config) logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	res := responsePayload{
		Message: "Logged out successfully",
		Error:   "",
		Data:    nil,
	}
	c.JSON(http.StatusOK, res)
}

func (app *Config) allUsers(c *gin.Context) {
	userId, _ := c.Get("userId")

	user, err := app.Models.User.GetByID(userId.(string))

	if err != nil {
		res := responsePayload{
			Message: "User does not exist",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	users, err := user.GetAllOtherUsersInOrg()

	if err != nil {
		res := responsePayload{
			Message: "Failed to get all users",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := responsePayload{
		Message: "Successfully get all other users in organization",
		Error:   "",
		Data: map[string]any{
			"users": users,
		},
	}
	c.JSON(http.StatusOK, res)
}

func (app *Config) addUser(c *gin.Context) {
	currentUserId, _ := c.Get("userId")

	currentUser, err := app.Models.User.GetByID(currentUserId.(string))

	if err != nil {
		res := responsePayload{
			Message: "User does not exist",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if currentUser.Role != "admin" {
		res := responsePayload{
			Message: "Not Authorized",
			Error:   "not authorized",
			Data:    nil,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	var reqPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err = c.Bind(&reqPayload)

	if err != nil {
		res := responsePayload{
			Message: "Error reading request body",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	username := reqPayload.Username
	password := reqPayload.Password

	if username == "" || password == "" {
		res := responsePayload{
			Message: "Missing Username or Password in request",
			Error:   "missing username or password in request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	userToAdd := data.User{
		Username:       username,
		Password:       password,
		OrganizationID: currentUser.OrganizationID,
		Role:           "member",
	}

	err = app.Models.User.Insert(userToAdd)

	if err != nil {
		res := responsePayload{
			Message: "Failed to add new user",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := responsePayload{
		Message: "Successfully add new user",
		Error:   "",
		Data: map[string]any{
			"user": userToAdd,
		},
	}
	c.JSON(http.StatusOK, res)
}

func (app *Config) deleteUser(c *gin.Context) {
	currentUserId, _ := c.Get("userId")

	currentUser, err := app.Models.User.GetByID(currentUserId.(string))

	if err != nil {
		res := responsePayload{
			Message: "User does not exist",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if currentUser.Role != "admin" {
		res := responsePayload{
			Message: "Not Authorized",
			Error:   "not authorized",
			Data:    nil,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	var reqPayload struct {
		Username string `json:"username"`
	}

	err = c.Bind(&reqPayload)

	if err != nil {
		res := responsePayload{
			Message: "Error reading request body",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	username := reqPayload.Username

	if username == "" {
		res := responsePayload{
			Message: "Missing Username in request",
			Error:   "missing username in request",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err = app.Models.User.Delete(username)

	if err != nil {
		res := responsePayload{
			Message: "Failed to delete user",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := responsePayload{
		Message: "Successfully delete user from organization",
		Error:   "",
		Data:    nil,
	}
	c.JSON(http.StatusOK, res)
}

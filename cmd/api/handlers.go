package main

import (
	"houseware---backend-engineering-octernship-KunalSin9h/data"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func sendResponse(message, err string, data map[string]any, c *gin.Context, code int) {
	var sendResponse struct {
		Message string         `json:"message"`
		Error   string         `json:"error"`
		Data    map[string]any `json:"data,omitempty"`
	}
	sendResponse.Message = message
	sendResponse.Error = err
	sendResponse.Data = data

	c.JSON(code, sendResponse)
}

/*
Login is a handler that takes the username and password from the request body and checks if the user exists.
If user exist in the data base then it create a JWT token and set it in the cookie.
*/
func (app *Config) login(c *gin.Context) {

	var reqPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := c.Bind(&reqPayload)

	if err != nil {
		sendResponse("Error reading request body", err.Error(), nil, c, http.StatusInternalServerError)
		return
	}

	username := reqPayload.Username
	password := reqPayload.Password

	if username == "" || password == "" {
		sendResponse("Missing Username or Password in request", "missing username or password in request", nil, c, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByUsername(username)

	if err != nil {
		sendResponse("Error while getting user", "error while getting user", nil, c, http.StatusInternalServerError)
		return
	}

	if user.ID == "" {
		// User Does not exist
		sendResponse("Invalid username of password", "invalid username or password", nil, c, http.StatusBadRequest)
		return
	}

	isPasswordMatched, err := user.PasswordMatch(password)

	if err != nil {
		sendResponse("Error while verifying password", err.Error(), nil, c, http.StatusInternalServerError)
		return
	}

	if !isPasswordMatched {
		// invalid password
		sendResponse("Invalid username or password", "invalid username or password", nil, c, http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(JWT_SECRET))

	if err != nil {
		sendResponse("Failed to create JWT Token", err.Error(), nil, c, http.StatusInternalServerError)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	sendResponse("User Signed in", "", map[string]any{
		"user": user,
	}, c, http.StatusOK)
}

/*
Logout is a handler that takes the JWT token from the cookie and set it to empty string.
*/
func (app *Config) logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	sendResponse("Logged out successfully", "", nil, c, http.StatusOK)
}

/*
AllUsers is a handler that return all other users in the organization.
*/
func (app *Config) allUsers(c *gin.Context) {
	userId, _ := c.Get("userId")

	user, err := app.Models.User.GetByID(userId.(string))

	if err != nil {
		sendResponse("User does not exist", err.Error(), nil, c, http.StatusBadRequest)
		return
	}

	users, err := user.GetAllOtherUsersInOrg()

	if err != nil {
		sendResponse("Failed to get all users", err.Error(), nil, c, http.StatusInternalServerError)
		return
	}

	sendResponse("Successfully get all other users in organization", "", map[string]any{
		"users": users,
	}, c, http.StatusOK)
}

/*
AddUser is a handler that takes the username and password from the request body and add a new user in the organization.
It can only be called by an admin.
*/
func (app *Config) addUser(c *gin.Context) {
	currentUserId, _ := c.Get("userId")

	currentUser, err := app.Models.User.GetByID(currentUserId.(string))

	if err != nil {
		sendResponse("User does not exist", err.Error(), nil, c, http.StatusBadRequest)
		return
	}

	if currentUser.Role != "admin" {
		sendResponse("Not Authorized", "not authorized", nil, c, http.StatusUnauthorized)
		return
	}

	var reqPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err = c.Bind(&reqPayload)

	if err != nil {
		sendResponse("Error reading request body", err.Error(), nil, c, http.StatusInternalServerError)
		return
	}

	username := reqPayload.Username
	password := reqPayload.Password

	if username == "" || password == "" {
		sendResponse("Missing Username or Password in request", "missing username or password in request", nil, c, http.StatusBadRequest)
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
		sendResponse("Failed to add new user", err.Error(), nil, c, http.StatusInternalServerError)
		return
	}

	sendResponse("Successfully add new user", "", map[string]any{
		"user": userToAdd,
	}, c, http.StatusOK)
}

/*
DeleteUser is a handler that takes the username from the request body and delete the user from the organization.
It can only be called by an admin.
*/
func (app *Config) deleteUser(c *gin.Context) {
	currentUserId, _ := c.Get("userId")

	currentUser, err := app.Models.User.GetByID(currentUserId.(string))

	if err != nil {
		sendResponse("Failed to get user", err.Error(), nil, c, http.StatusBadRequest)
		return
	}

	if currentUser.Role != "admin" {
		sendResponse("Not Authorized", "not authorized", nil, c, http.StatusUnauthorized)
		return
	}

	var reqPayload struct {
		Username string `json:"username"`
	}

	err = c.Bind(&reqPayload)

	if err != nil {
		sendResponse("Error reading request body", err.Error(), nil, c, http.StatusInternalServerError)
		return
	}

	username := reqPayload.Username

	if username == "" {
		sendResponse("Missing username in request", "missing username in request", nil, c, http.StatusBadRequest)
		return
	}

	userToDelete, err := app.Models.User.GetByUsername(username)

	if err != nil {
		sendResponse("Failed to delete user", err.Error(), nil, c, http.StatusInternalServerError)
		return
	}

	if userToDelete.ID == "" || currentUser.OrganizationID != userToDelete.OrganizationID {
		sendResponse("Not Authorized", "not authorized", nil, c, http.StatusUnauthorized)
		return
	}

	err = userToDelete.Delete()

	if err != nil {
		sendResponse("Failed to delete user", err.Error(), nil, c, http.StatusBadRequest)
		return
	}

	sendResponse("Successfully delete user from organization", "", nil, c, http.StatusOK)
}

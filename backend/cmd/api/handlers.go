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

type SignUpRequestPayload struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Role           string `json:"role"`
	OrganizationId string `json:"organization_id"`
}

type loginRequestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *Config) signUp(c *gin.Context) {

	var reqPayload SignUpRequestPayload

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

	user := data.User{
		FirstName:      reqPayload.FirstName,
		LastName:       reqPayload.LastName,
		Email:          reqPayload.Email,
		Username:       reqPayload.Username,
		Password:       reqPayload.Password,
		Role:           reqPayload.Role,
		OrganizationID: reqPayload.OrganizationId,
	}

	if user.FirstName == "" ||
		user.LastName == "" ||
		user.Email == "" ||
		user.Username == "" ||
		user.Password == "" ||
		user.Role == "" ||
		user.OrganizationID == "" {
		res := responsePayload{
			Message: "Missing user details",
			Error:   "missing user details",
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err = app.Models.User.Insert(user)

	if err != nil {
		res := responsePayload{
			Message: "Failed to insert user",
			Error:   err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := responsePayload{
		Message: "User signed up",
		Error:   "",
		Data:    nil,
	}

	c.JSON(http.StatusOK, res)
}

func (app *Config) login(c *gin.Context) {

	var reqPayload loginRequestPayload
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

	user, err := app.Models.User.GetUserByUsername(username)

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

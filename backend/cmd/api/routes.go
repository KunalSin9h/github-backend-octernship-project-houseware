package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/*
routes returns all the routes for the application.
*/
func (app *Config) routes() *gin.Engine {

	router := gin.New()

	/*
		Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
		By default, gin.DefaultWriter = os.Stdout.
	*/
	router.Use(gin.Logger())

	/*
		Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
	*/
	router.Use(gin.Recovery())

	/*
		Enabling Cors
	*/
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	/* Grouping by version */
	v1 := router.Group("/v1")

	/* Registering Routes */
	v1.POST("/login", app.login)   // User Login
	v1.POST("/logout", app.logout) // User Logout

	// Admin User adds a new User account(by providing the username & password)
	v1.POST("/add", app.AuthorizationMiddleware, app.addUser)

	// Admin User deletes an existing User account from their organization
	v1.DELETE("/delete", app.AuthorizationMiddleware, app.deleteUser)

	//List all Users in their organization
	v1.GET("/users", app.AuthorizationMiddleware, app.allUsers)

	return router
}

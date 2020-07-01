package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func signupForm(c echo.Context) error {
	// Get name
	firstName := c.FormValue("firstname")
	lastName := c.FormValue("lastname")
	email := c.FormValue("email")
	username := c.FormValue("username")

	fmt.Println("First name:", firstName)
	fmt.Println("Last name:", lastName)
	fmt.Println("Email address:", email)
	fmt.Println("Username:", username)

	return c.String(http.StatusOK, "First name: "+firstName+"\nLast name: "+lastName+"\nEmail: "+email+"\nUsername: "+username)
}

func deleteForm(c echo.Context) error {
	// Get name
	firstName := c.FormValue("firstname")
	lastName := c.FormValue("lastname")
	email := c.FormValue("email")

	fmt.Println("First name:", firstName)
	fmt.Println("Last name:", lastName)
	fmt.Println("Email address:", email)

	return c.String(http.StatusOK, "First name: "+firstName+"\nLast name: "+lastName+"\nEmail: "+email)
}

func main() {
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderAuthorization, echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	}))

	e.POST("/api/signup", signupForm)
	e.POST("/api/delete", deleteForm)
	e.Logger.Fatal(e.Start(":1323"))
}

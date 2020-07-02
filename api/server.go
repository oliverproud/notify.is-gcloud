package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"notify.is-go/database"
)

func signupForm(c echo.Context) error {
	// Get name
	firstName := strings.Title(strings.ToLower(c.FormValue("firstname")))
	lastName := strings.Title(strings.ToLower(c.FormValue("lastname")))
	email := strings.ToLower(c.FormValue("email"))
	username := strings.ToLower(c.FormValue("username"))

	fmt.Println("First name:", firstName)
	fmt.Println("Last name:", lastName)
	fmt.Println("Email address:", email)
	fmt.Println("Username:", username)

	result, err := models.InsertUser(firstName, lastName, email, username)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	return c.String(http.StatusOK, "First name: "+firstName+"\nLast name: "+lastName+"\nEmail: "+email+"\nUsername: "+username)
}

func deleteForm(c echo.Context) error {
	// Get name
	firstName := strings.Title(strings.ToLower(c.FormValue("firstname")))
	lastName := strings.Title(strings.ToLower(c.FormValue("lastname")))
	email := strings.ToLower(c.FormValue("email"))

	fmt.Println("First name:", firstName)
	fmt.Println("Last name:", lastName)
	fmt.Println("Email address:", email)

	result, err := models.DeleteUser(firstName, lastName, email)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	return c.String(http.StatusOK, "First name: "+firstName+"\nLast name: "+lastName+"\nEmail: "+email)
}

const (
	host   = "localhost"
	port   = 5432
	user   = "oliverproud"
	dbname = "notify"
)

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

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	models.InitDB(psqlInfo)

	defer models.CloseDB()

	e.POST("/api/signup", signupForm)
	e.POST("/api/delete", deleteForm)
	e.Logger.Fatal(e.Start(":1323"))
}

package main

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Group level middleware
	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		if username == "joe" && password == "secret" {
			ctx.Set("username", username)
			return true, nil
		}
		return false, nil
	}))

	// Routes
	e.GET("/", func(c echo.Context) error {
		data := map[string]string{
			"message":     "Hello, World",
			"description": "This is the hello world description",
		}
		return c.JSONPretty(http.StatusOK, data, " ")
	})

	e.GET("/users/:id", getUser)
	e.GET("/show", show)
	e.POST("/save", save)
	e.POST("/multipartsave", multiPartSave)

	adminDashboardHandler := func(c echo.Context) error {
		username := c.Get("username").(string)
	
		// Return a JSON response
		response := map[string]interface{}{
			"message": "Welcome to the admin dashboard, " + username + "!",

		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	
	g.GET("/dashboard", adminDashboardHandler)

	// Start the server
	e.Logger.Fatal(e.Start(":1323"))
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func show(c echo.Context) error {
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

func save(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, name+" "+email)
}

func multiPartSave(c echo.Context) error {
	name := c.FormValue("name")
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(avatar.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<b>Thank you! "+name+"</b>")
}

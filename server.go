package main

import (
	"testing/auth"
	"testing/helpers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	helpers.ConnectToDB();
	e.Use(middleware.Logger());
	e.Use(middleware.Recover());

	e.POST("/auth/login", auth.LoginHandler);
	e.POST("/auth/register", auth.SignupHandler);

	e.Logger.Fatal(e.Start(":1323"))
	helpers.DisconnectFromDB();
}

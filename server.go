package main

import (
	"github.com/woaitsAryan/portfolio-website/auth"
	"github.com/woaitsAryan/portfolio-website/helpers"
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

	e.Logger.Fatal(e.Start(":" + helpers.PORT))
	helpers.DisconnectFromDB();
}

package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/woaitsAryan/portfolio-website/auth"
	"github.com/woaitsAryan/portfolio-website/datafetch"
	"github.com/woaitsAryan/portfolio-website/helpers"
	"github.com/woaitsAryan/portfolio-website/middlewares"
)

func main() {
	e := echo.New()
	helpers.ConnectToDB();
	e.Use(middleware.Logger());
	e.Use(middleware.Recover());
	e.POST("/auth/login", auth.LoginHandler);
	e.POST("/auth/register", auth.SignupHandler);
	e.POST("/user/write", middlewares.ProtectJWT(datafetch.WriteData));
	e.POST("/user/get", middlewares.ProtectJWT(datafetch.GetData));
	e.Logger.Fatal(e.Start(":" + helpers.PORT))
	helpers.DisconnectFromDB();
}

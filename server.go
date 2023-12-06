package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/woaitsAryan/portfolio-website/auth"
	"github.com/woaitsAryan/portfolio-website/CRUD"
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
	e.POST("/user/create", middlewares.ProtectJWT(CRUD.CreateData));
	e.POST("/user/read", middlewares.ProtectJWT(CRUD.ReadData));
	e.POST("/user/update", middlewares.ProtectJWT(CRUD.UpdateData));
	e.POST("/user/delete", middlewares.ProtectJWT(CRUD.DeleteData));
	e.Logger.Fatal(e.Start(":" + helpers.PORT))
	helpers.DisconnectFromDB();
}

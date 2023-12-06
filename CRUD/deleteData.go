package CRUD

import (
	"context"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/woaitsAryan/portfolio-website/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteData(c echo.Context) error {
	result := c.Get("result")

	resultData, _ := result.(map[string]interface{})
	storedEmail, _ := resultData["email"].(string)

	filter := bson.M{"email": storedEmail }
	update := bson.M{"$set": bson.M{"data": ""}} 

	_ = helpers.Usercollection.FindOneAndUpdate(context.Background(), filter, update)
	return c.JSON(
		http.StatusFound, map[string]string{
			"message": "Data deleted!",
		},
	)
}

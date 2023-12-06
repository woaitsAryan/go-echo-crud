package datafetch

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/woaitsAryan/portfolio-website/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

type reqbody struct{
	Data string `json:"data"`
}

func WriteData(c echo.Context) error {
	var body reqbody
	if err := c.Bind(&body); err != nil || body.Data == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}
	data := body.Data
	result := c.Get("result")
	resultData, _ := result.(map[string]interface{})
	storedEmail, _ := resultData["email"].(string)

	if data == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}
	filter := bson.M{ "email": storedEmail }
	update := bson.M{"$set": bson.M{"data": data}} 

	_ = helpers.Usercollection.FindOneAndUpdate(context.Background(), filter, update)
	
	return c.JSON(
		http.StatusCreated, map[string]string{
			"message": "Data successfully added",
		},
	)
}

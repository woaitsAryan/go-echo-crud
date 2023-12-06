package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/woaitsAryan/portfolio-website/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProtectJWT(next echo.HandlerFunc) echo.HandlerFunc{
	return func(c echo.Context) error{
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == ""{
			return c.JSON(
				http.StatusUnauthorized,
				map[string]any{"message": "Unauthorized access"},
			)
		}
		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			return []byte(helpers.JWT_KEY), nil
		})
		if (err != nil || !token.Valid){
			return c.JSON(
				http.StatusUnauthorized,
				map[string]any{"message": "Unauthorized access"},
			)
		}
		var result bson.M
		claims, ok := token.Claims.(jwt.MapClaims)
		if(!ok){
			return c.JSON(
				http.StatusUnauthorized,
				map[string]any{"message": "Unauthorized access"},
			)
		}
		id, err := primitive.ObjectIDFromHex(claims["userID"].(string))

		if(err != nil){
			return c.JSON(
				http.StatusUnauthorized,
				map[string]any{"message": "Unauthorized access"},
			)
		}

		filter := bson.M{"_id": id}

		err = helpers.Usercollection.FindOne(context.Background(), filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "User with this email does not exist",
				})
			}
			log.Fatal(err)
		}
		resultMap := map[string]interface{}(result)
		c.Set("result", resultMap)
		return next(c);
	}
}
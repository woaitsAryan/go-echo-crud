package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing/helpers"
	"testing/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil || user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	existingUser := bson.M{"email": user.Email}
	count, err := helpers.Usercollection.CountDocuments(context.Background(), existingUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Error in fetching",
		})
	}

	if count > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "User with this email already exists",
		})
	}
	fmt.Println("Reached here")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "internal server error",
		})
	}
	newUser := bson.M{
		"email":    user.Email,
		"password": string(hashedPassword),
	}

	mongouser, err := helpers.Usercollection.InsertOne(context.Background(), newUser)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Error inserting user",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": mongouser.InsertedID.(primitive.ObjectID).String(),
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Signing failure",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Signed Up successfully!",
		"token":   tokenString,
	})
}

func LoginHandler(c echo.Context) error {

	user := new(models.User)
	if err := c.Bind(user); err != nil || user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	var result bson.M
	filter := bson.M{"email": user.Email}
	err := helpers.Usercollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "User with this email does not exist",
			})
		}
		log.Fatal(err)
	}
	storedPassword := result["password"].(string)
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid credentials",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": result["_id"].(primitive.ObjectID).String(),
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Signing failure",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login successful!",
		"token":   tokenString,
	})

}

package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/woaitsAryan/portfolio-website/helpers"
	"github.com/woaitsAryan/portfolio-website/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "internal server error",
		})
	}
	newUser := bson.M{
		"email":    user.Email,
		"password": string(hashedPassword),
		"data": "",
	}

	mongouser, err := helpers.Usercollection.InsertOne(context.Background(), newUser)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Error inserting user",
		})
	}

	userClaim := jwt.MapClaims{
		"userID": mongouser.InsertedID.(primitive.ObjectID).Hex(),
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)

	tokenString, err := token.SignedString([]byte(helpers.JWT_KEY))

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
		"userID": result["_id"].(primitive.ObjectID).Hex(),
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(helpers.JWT_KEY))

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

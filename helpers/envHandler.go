package helpers

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

var (
	PORT      string = os.Getenv("PORT")
	MONGO_URL string = os.Getenv("MONGO_URL")
	JWT_KEY   string = os.Getenv("JWT_KEY")
)

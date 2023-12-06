package models

type User struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Data     string `json:"data" bson:"data"`
}

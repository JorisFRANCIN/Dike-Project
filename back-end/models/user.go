package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Password      *string            `json:"password" form:"password" validate:"required,min=6"`
	Email         *string            `json:"email" form:"email" validate:"email,required"`
	Token         *string            `json:"token" form:"token"`
	Refresh_Token *string            `json:"refresh_token" form:"refresh_token"`
	Created_at    time.Time          `json:"created_at" form:"created_at"`
	Updated_at    time.Time          `json:"updated_at" form:"updated_at"`
	User_id       string             `json:"user_id" form:"user_id"`
	Username      *string            `json:"username" form:"username"`
	Services      []ServiceAccessToken `json:"services"`
}

type UserRegister struct {
	Password      *string            `json:"password" form:"password" validate:"required,min=6"`
	Email         *string            `json:"email" form:"email" validate:"email,required"`
	Username      *string            `json:"username" form:"username"`
}

type UserLogin struct {
	Email         *string            `json:"email" form:"email" validate:"email,required"`
	Password      *string            `json:"password" form:"password" validate:"required,min=6"`
}

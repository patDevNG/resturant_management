package models

import(
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID                 `bson:"_id"`
	First_name    string                             `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     string                             `json:"last_name" validate:"required"`
	Password      string                             `json:"password" validate:"required"`
	Email         string                             `json:"email" validate:"required"`
	Avartar       *string                            `json:"avatar"`  
	Phone         *string                            `json:"phone" validate:"required"`
	Token         *string                            `json:"token"`
	Refresh_Token *string                            `json:"referesh_token"` 
	User_id       string                             `json:"user_id"`
	Created_at    time.Time                          `json:"createdAt"`
	Updated_at    time.Time                          `json:"updated_at"`
}
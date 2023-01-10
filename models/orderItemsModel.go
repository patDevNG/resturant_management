package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItems struct{
	ID              primitive.ObjectID            `bson:"_id"`
	Order_id        string                       `json:"order_id" validate:"required"`
	Order_items_id  string                        `json:"order_item_id"`
	Quantity        *int                           `json:"quantity" validate:"required"`
	Unit_price      *float64                      `json:"unit_price" validate:"required"`
	Food_id         *string                       `json:"food_id"`
	Created_at      time.Time                     `json:"created_at"`
	Updated_at      time.Time                     `json:"updated_at"`
}
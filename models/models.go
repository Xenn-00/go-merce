package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	FirstName       *string            `bson:"first_name" json:"firstName" validate:"required,min=2,max=100"`
	LastName        *string            `bson:"last_name" json:"lastName" validate:"required,min=2,max=100"`
	Password        *string            `bson:"password" json:"password" validate:"required,min=8"`
	Email           *string            `bson:"email" json:"email" validate:"email,required"`
	Phone           *string            `bson:"phone" json:"phone" validate:"required"`
	Token           *string            `bson:"token" json:"token"`
	Refresh_Token   *string            `bson:"refresh_token" json:"refresh_token"`
	Created_At      time.Time          `bson:"created_at" json:"created_at"`
	Updated_At      time.Time          `bson:"updated_at" json:"updated_at"`
	User_ID         string             `bson:"userID" json:"userID"`
	UserCart        []ProductUser      `bson:"user_cart" json:"user_cart"`
	Address_Details []Address          `bson:"address_details" json:"address_details"`
	Order_Status    []Order            `bson:"order_status" json:"order_status"`
}

type Product struct {
	Product_ID   primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Product_Name *string            `bson:"product_name" json:"product_name"`
	Price        *int               `bson:"price" json:"price"`
	Rating       *int               `bson:"rating" json:"rating"`
	Image        *string            `bson:"image" json:"image"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Product_Name *string            `bson:"product_name" json:"product_name"`
	Price        int                `bson:"price" json:"price"`
	Rating       *uint              `bson:"rating" json:"rating"`
	Image        *string            `bson:"image" json:"image"`
}

type Address struct {
	Address_ID primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	House      *string            `bson:"house" json:"house"`
	Street     *string            `bson:"street" json:"street"`
	City       *string            `bson:"city" json:"city"`
	Province   *string            `bson:"province" json:"province"`
	Postcode   *string            `bson:"postcode" json:"postcode"`
}

type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Order_Cart     []ProductUser      `bson:"order_cart" json:"order_cart"`
	Ordered_At     time.Time          `bson:"ordered_at" json:"ordered_at"`
	Price          int                `bson:"price" json:"price"`
	Discount       *int               `bson:"discount" json:"discount"`
	Payment_method Payment            `bson:"payment_method" json:"payment_method"`
}

type Payment struct {
	Digital bool `bson:"digital" json:"digital"`
	COD     bool `bson:"cod" json:"cod"`
}

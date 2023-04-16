package database

import (
	"context"
	"errors"
	"log"

	"github.com/Xenn-00/go-merce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("sorry, we can't find the products")
	ErrCantDecodeProducts = errors.New("sorry, we can't find the products")
	ErrUserIdIsNotValid   = errors.New("sorry, this user is not valid")
	ErrCantUpdateUser     = errors.New("can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("can't remove this item from the cart")
	ErrCantGetItems       = errors.New("unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("can't update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, prodID primitive.ObjectID, userId string) error {
	// check product from db first
	searchFromDB, err := prodCollection.Find(ctx, bson.M{"_id": prodID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	err = searchFromDB.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{
		primitive.E{
			Key:   "_id",
			Value: id,
		},
	}
	update := bson.D{
		{
			Key: "$push",
			Value: bson.D{
				primitive.E{
					Key: "user_cart",
					Value: bson.D{
						{
							Key:   "$each",
							Value: productCart,
						},
					},
				},
			},
		},
	}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}

	return err
}

func RemoveCartItem()

func BuyItemFromCart()

func InstantBuy()

package database

import (
	"context"
	"errors"
	"log"
	"time"

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

func RemoveCartItem(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productID primitive.ObjectID, userId string) error {
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
	update := bson.M{
		"$pull": bson.M{
			"user_cart": bson.M{
				"_id": productID,
			},
		},
	}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var getCartItems models.User
	var orderCart models.Order
	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Ordered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_method.COD = true
	unwind := bson.D{
		{
			Key: "$unwind",
			Value: bson.D{
				primitive.E{
					Key:   "path",
					Value: "$user_cart",
				},
			},
		},
	}
	grouping := bson.D{{
		Key: "$group",
		Value: bson.D{
			primitive.E{
				Key:   "_id",
				Value: "$_id",
			},
			{
				Key: "total",
				Value: bson.D{
					primitive.E{
						Key:   "$sum",
						Value: "$user_cart.price",
					},
				},
			},
		},
	}}
	currentResult, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}
	var getUserCart []bson.M
	if err = currentResult.All(ctx, &getUserCart); err != nil {
		panic(err)
	}
	var totalPrice int32
	for _, user_item := range getUserCart {
		price := user_item["total"]
		totalPrice = price.(int32)
	}
	orderCart.Price = int(totalPrice)
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
					Key:   "orders",
					Value: orderCart,
				},
			},
		},
	}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}
	secondFilter := bson.D{
		primitive.E{
			Key:   "_id",
			Value: id,
		},
	}
	secondUpdate := bson.M{
		"$push": bson.M{
			"orders.$[].order_list": bson.M{
				"$each": getCartItems.UserCart,
			},
		},
	}
	_, err = userCollection.UpdateOne(ctx, secondFilter, secondUpdate)
	if err != nil {
		log.Println(err)

	}
	emptyUserCart := make([]models.ProductUser, 0)
	filtered := bson.D{
		primitive.E{
			Key:   "_id",
			Value: id,
		},
	}
	updated := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				primitive.E{
					Key:   "user_cart",
					Value: emptyUserCart,
				},
			},
		},
	}
	_, err = userCollection.UpdateOne(ctx, filtered, updated)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuy(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var productDetails models.ProductUser
	var ordersDetail models.Order
	ordersDetail.Order_ID = primitive.NewObjectID()
	ordersDetail.Ordered_At = time.Now()
	ordersDetail.Order_Cart = make([]models.ProductUser, 0)
	ordersDetail.Payment_method.COD = true
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
	}
	ordersDetail.Price = productDetails.Price
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
					Key:   "orders",
					Value: ordersDetail,
				},
			},
		},
	}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	secondFilter := bson.D{
		primitive.E{
			Key:   "_id",
			Value: id,
		},
	}
	secondUpdate := bson.M{
		"$push": bson.M{
			"orders.$[].order_list": productDetails,
		},
	}
	_, err = userCollection.UpdateOne(ctx, secondFilter, secondUpdate)
	if err != nil {
		log.Println(err)
	}
	return nil
}

package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Xenn-00/go-merce/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	ProductCollection *mongo.Collection
	UserCollection    *mongo.Collection
}

func NewApplication(ProductCollection *mongo.Collection, UserCollection *mongo.Collection) *Application {
	return &Application{
		ProductCollection: ProductCollection,
		UserCollection:    UserCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       errors.New("product id is empty"),
			})
			return
		}

		userQueryID := c.Query("User_ID")
		if userQueryID == "" {
			log.Println("user id is empty")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       errors.New("user id is empty"),
			})
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.ProductCollection, app.UserCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
		}
		c.IndentedJSON(http.StatusOK, "Successfully add to cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       errors.New("product id is empty"),
			})
			return
		}

		userQueryID := c.Query("User_ID")
		if userQueryID == "" {
			log.Println("user id is empty")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       errors.New("user id is empty"),
			})
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.ProductCollection, app.UserCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
		}
		c.IndentedJSON(http.StatusOK, "Successfully remove item from cart")
	}
}

func GetItemFromCart() gin.HandlerFunc

func BuyFromCart() gin.HandlerFunc

func InstantBuy() gin.HandlerFunc

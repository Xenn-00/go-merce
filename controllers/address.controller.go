package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Xenn-00/go-merce/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc

func EditAddress() gin.HandlerFunc

func EditWorkAddress() gin.HandlerFunc

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		Id := c.Query("id")

		if Id == "" {
			c.Header("Content-type", "application/json")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       "Invalid search index",
			})
			return
		}
		// technically it's not remove the value, it's just update the value to 0 / empty
		addresses := make([]models.Address, 0)
		user_id, err := primitive.ObjectIDFromHex(Id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"status_code": http.StatusInternalServerError,
				"error":       err.Error(),
			})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "Wrong Command")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully deleted")
	}
}

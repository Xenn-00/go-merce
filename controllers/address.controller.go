package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Xenn-00/go-merce/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Invalid id",
			})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"status_code": http.StatusInternalServerError,
				"error":       err.Error(),
			})
			return
		}

		var addresses models.Address
		addresses.Address_ID = primitive.NewObjectID()
		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, gin.H{
				"error": err.Error(),
			})
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		filter_match := bson.D{
			{Key: "$match",
				Value: bson.D{
					primitive.E{
						Key:   "_id",
						Value: address,
					},
				},
			},
		}
		unwind := bson.D{
			{Key: "$unwind",
				Value: bson.D{
					primitive.E{
						Key:   "path",
						Value: "$address",
					},
				},
			},
		}
		group := bson.D{{
			Key: "$group",
			Value: bson.D{
				primitive.E{
					Key:   "_id",
					Value: "$address_id",
				},
				{
					Key: "count",
					Value: bson.D{
						primitive.E{
							Key:   "$sum",
							Value: 1,
						},
					},
				},
			},
		}}

		cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, group})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		var addressInfo []bson.M
		if err := cursor.All(ctx, &addressInfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{
				primitive.E{
					Key:   "_id",
					Value: address,
				},
			}
			update := bson.D{
				{
					Key: "$push",
					Value: bson.D{
						primitive.E{
							Key:   "address",
							Value: addresses,
						},
					},
				},
			}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			c.IndentedJSON(http.StatusBadRequest, "Not Allowed")
		}
		defer cancel()
		ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		Id := c.Query("id")
		if Id == "" {
			c.Header("Content-type", "application/json")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       "Invalid",
			})
			return
		}
		user_id, err := primitive.ObjectIDFromHex(Id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"status_code": http.StatusInternalServerError,
				"error":       err.Error(),
			})
			return
		}
		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       err.Error(),
			})
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{
			primitive.E{
				Key:   "_id",
				Value: user_id,
			},
		}
		update := bson.D{{
			Key: "$set",
			Value: bson.D{
				primitive.E{
					Key:   "address.0.house",
					Value: editAddress.House,
				},
				{
					Key:   "address.0.street",
					Value: editAddress.Street,
				},
				{
					Key:   "address.0.city",
					Value: editAddress.City,
				},
				{
					Key:   "address.0.province",
					Value: editAddress.Province,
				},
				{
					Key:   "address.0.postcode",
					Value: editAddress.Postcode,
				},
			},
		}}
		if _, err := UserCollection.UpdateOne(ctx, filter, update); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully update house address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		Id := c.Query("id")
		if Id == "" {
			c.Header("Content-type", "application/json")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       "Invalid",
			})
			return
		}
		user_id, err := primitive.ObjectIDFromHex(Id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"status_code": http.StatusInternalServerError,
				"error":       err.Error(),
			})
			return
		}
		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       err.Error(),
			})
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{
			primitive.E{
				Key:   "_id",
				Value: user_id,
			},
		}
		update := bson.D{{
			Key: "$set",
			Value: bson.D{
				primitive.E{
					Key:   "address.1.house",
					Value: editAddress.House,
				},
				{
					Key:   "address.1.street",
					Value: editAddress.Street,
				},
				{
					Key:   "address.1.city",
					Value: editAddress.City,
				},
				{
					Key:   "address.1.province",
					Value: editAddress.Province,
				},
				{
					Key:   "address.1.postcode",
					Value: editAddress.Postcode,
				},
			},
		}}
		if _, err := UserCollection.UpdateOne(ctx, filter, update); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully update work address")
	}
}

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

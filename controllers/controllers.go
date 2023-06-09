package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Xenn-00/go-merce/database"
	"github.com/Xenn-00/go-merce/models"
	"github.com/Xenn-00/go-merce/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err.Error())
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	if err != nil {
		return false, "incorrect email or password"
	}
	return true, ""
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       err.Error(),
			})

			return
		}
		validate := validator.New()
		err := validate.Struct(user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       err,
			})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": http.StatusInternalServerError,
				"error":       err.Error(),
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       "user already exist",
			})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": http.StatusInternalServerError,
				"error":       err.Error(),
			})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       "phone number is already in use",
			})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refresh_token, _ := tokens.TokenGenerator(*user.Email, *user.FirstName, *user.LastName, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refresh_token
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": http.StatusInternalServerError,
				"error":       "can't create user",
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusCreated, gin.H{
			"status_code": http.StatusCreated,
			"message":     "successfully signed in",
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       err.Error(),
			})
			return
		}
		var foundUser models.User
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status_code": http.StatusInternalServerError,
				"error":       "incorrect email or password",
			})
			return
		}

		isPasswordValid, message := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !isPasswordValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"error":       "incorrect email or password",
			})
			fmt.Println(message)
			return
		}
		token, refresh_token, _ := tokens.TokenGenerator(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *&foundUser.User_ID)
		defer cancel()

		tokens.UpdateAllToken(token, refresh_token, foundUser.User_ID)
		c.JSON(http.StatusFound, gin.H{
			"status_code": http.StatusFound,
			"message":     "successfully login",
		})
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()
		if err := c.ShouldBindJSON(&products);err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error" : err.Error(),
			})
			return
		}
		products.Product_ID = primitive.NewObjectID()
		_, err := ProductCollection.InsertOne(ctx, products)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error" : err.Error(),
			})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"status_code" : http.StatusOK,
			"message" : "Successfully add product",
		})
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something when wrong, please try again")
			return
		}

		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err.Error())
			c.IndentedJSON(http.StatusBadRequest, "invalid")
			return
		}

		defer cancel()
		c.IndentedJSON(http.StatusOK, productList)
	}
}

func SearchProductbyQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProducts []models.Product
		queryParam := c.Query("name")

		// check if it's empty
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{
				"status_code": http.StatusNotFound,
				"error":       "Invalid search index",
			})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchQueryDB, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "something went wrong while fetching the data")
			return
		}
		err = searchQueryDB.All(ctx, &searchProducts)
		if err != nil {
			log.Println(err.Error())
			c.IndentedJSON(http.StatusNotFound, "invalid")
			return
		}
		defer searchQueryDB.Close(ctx)

		if err := searchQueryDB.Err(); err != nil {
			log.Println(err.Error())
			c.IndentedJSON(http.StatusBadRequest, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(http.StatusOK, searchProducts)
	}
}

package controller

import (
	"context"
	"fmt"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food models.Food
		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the food item"})

		}
		c.JSON(http.StatusOK, food)
		defer cancle()
	}
}
func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		var food models.Food
		var menu models.Menu
		err := c.BindJSON(&food)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		validationErr := validate.Struct(food)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}
		defer cancle()
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		if err != nil {

			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now()).Format(time.RFC3339)
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now()).Format(time.RFC3339)
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num
		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := fmt.Sprintf("Food item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancle()
		c.JSON(http.StatusOK, result)

	}
}
func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func round(num float64) int {

}

func toFixed(num float64, precision int) float64 {

}

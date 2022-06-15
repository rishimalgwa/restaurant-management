package controller

import (
	"context"
	"golang-restaurant-management/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing the menu items"})
		}
		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)
	}
}
func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

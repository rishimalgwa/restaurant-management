package controller

import (
	"context"
	"fmt"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancle()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items"})

		}
		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders)

	}
}
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("order_id")
		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the order"})

		}
		c.JSON(http.StatusOK, order)
		defer cancle()
	}
}
func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var table models.Table
		var order models.Order
		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		err := c.BindJSON(&order)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(order)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		if order.Table_id != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)

			defer cancle()
			if err != nil {
				msg := fmt.Sprintf("table not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
		}
		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()
		result, insertErr := orderCollection.InsertOne(ctx, order)
		if insertErr != nil {
			msg := fmt.Sprintf("order was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancle()
		c.JSON(http.StatusOK, result)

	}
}
func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var table models.Table
		var order models.Order
		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		err := c.BindJSON(&order)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		orderID := c.Param("order_id")
		filter := bson.M{"order_id": orderID}
		var updateObj primitive.D
		if order.Table_id != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancle()
			if err != nil {
				MSG := fmt.Sprintf("order was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": MSG})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "table_id", Value: order.Table_id})
			order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{"updated_at", order.Updated_at})
			upsert := true
			opt := options.UpdateOptions{
				Upsert: &upsert,
			}
			result, err := orderCollection.UpdateOne(
				ctx,
				filter,
				bson.D{{"$set", updateObj}},
				&opt,
			)
			if err != nil {
				msg := fmt.Sprintf("order update failed")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			defer cancle()
			c.JSON(http.StatusOK, result)
		}
	}
}
func OrderItemOrderCreator(order models.Order) string {
	var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()
	orderCollection.InsertOne(ctx, order)
	defer cancle()
	return order.Order_id
}

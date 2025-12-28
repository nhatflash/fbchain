package database

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)


func ValidateRestaurantItemSchema(db *mongo.Database) error {
	collectionName := "restaurant_items"
	schema := bson.M{
		"bsonType": "object",
		"required": []string{"name", "price", "type", "status", "restaurantId"},
		"properties": bson.M{
			"name": bson.M{
				"bsonType": "string",
			},
			"price": bson.M{
				"bsonType": "decimal",
			},
			"type": bson.M{
				"bsonType": "string",
				"enum": []string{"FOOD", "BEVERAGE"},
			},
			"status": bson.M{
				"bsonType": "string",
				"enum": []string{"AVAILABLE", "UNAVAILABLE", "DELETED"},
			},
			"restaurantId": bson.M{
				"bsonType": "long",
			},
		},


	}
	opts := options.CreateCollection().SetValidator(bson.M{
		"$jsonSchema": schema,
	})
	err := db.CreateCollection(context.TODO(), collectionName, opts)
	if err != nil {
		if !mongo.IsDuplicateKeyError(err) && !strings.Contains(err.Error(), "already exists") {
			return err
		} 
	}
	return nil
}
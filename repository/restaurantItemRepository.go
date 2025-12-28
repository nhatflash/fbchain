package repository

import (
	"context"

	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RestaurantItemRepository struct {
	rItemColl 			*mongo.Collection
}


func NewRestaurantItemRepository(coll *mongo.Collection) *RestaurantItemRepository {
	return &RestaurantItemRepository{
		rItemColl: coll,
	}
}


func (rir *RestaurantItemRepository) AddNewRestaurantItem(ctx context.Context, rItem *model.RestaurantItem) (*model.RestaurantItem, error) {
	result, err := rir.rItemColl.InsertOne(ctx, rItem)
	if err != nil {
		return nil, err
	}
	var newItem *model.RestaurantItem
	newId := result.InsertedID
	newItem, err = rir.FindRestaurantItemById(ctx, newId.(string))
	if err != nil {
		return nil, err
	}
	return newItem, nil
}


func (rir *RestaurantItemRepository) FindRestaurantItemById(ctx context.Context, id string) (*model.RestaurantItem, error) {
	filter := bson.D{
		{Key: "_id", Value: id},
	}
	var item model.RestaurantItem
	err := rir.rItemColl.FindOne(ctx, filter, nil).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, appErr.NotFoundError("No restaurant item found.")
		}
		return nil, err
	}
	return &item, nil
}


func (rir *RestaurantItemRepository) FindAllRestaurantItems(ctx context.Context) ([]model.RestaurantItem, error) {
	cursor, err := rir.rItemColl.Find(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	var items []model.RestaurantItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (rir *RestaurantItemRepository) FindRestaurantItemsByRestaurantId(ctx context.Context, restaurantId int64) ([]model.RestaurantItem, error) {
	filter := bson.D{
		{Key: "restaurantId", Value: restaurantId},
	}
	cursor, err := rir.rItemColl.Find(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	var items []model.RestaurantItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}
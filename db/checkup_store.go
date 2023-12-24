package db

import (
	"context"
	"fmt"
	"os"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var chcoll = "checkups"

type CheckUpStore interface {
	GetCheckUps(context.Context, string) ([]*types.CheckUp, error)
	GetCheckByID(context.Context, string) (*types.CheckUp, error)
	InsertCheckUp(context.Context, *types.CheckUp) (*types.CheckUp, error)
	PerformCheckUp(context.Context, bson.D, bson.D) (*types.CheckUp, error)
}

type MongoCheckUpStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoCheckUpStore(client *mongo.Client) *MongoCheckUpStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoCheckUpStore{
		client: client,
		coll:   client.Database(dbname).Collection(chcoll),
	}
}

func (c *MongoCheckUpStore) PerformCheckUp(ctx context.Context, filter bson.D, update bson.D) (*types.CheckUp, error) {
	var cu *types.CheckUp
	if err := c.coll.FindOneAndUpdate(ctx, filter, update).Decode(&cu); err != nil {
		return nil, err
	}
	return cu, nil
}

func (c *MongoCheckUpStore) InsertCheckUp(ctx context.Context, ch *types.CheckUp) (*types.CheckUp, error) {

	checkup, err := c.coll.InsertOne(ctx, ch)
	if err != nil {
		return nil, err
	}
	ch.ID = checkup.InsertedID.(primitive.ObjectID)
	return ch, nil
}

func (c *MongoCheckUpStore) GetCheckUps(ctx context.Context, id string) ([]*types.CheckUp, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	cur, err := c.coll.Find(ctx, bson.M{"doctorId": oid})
	if err != nil {
		fmt.Println("error in GetCheckUps")
		return nil, err
	}
	var checkups []*types.CheckUp
	if err := cur.All(ctx, &checkups); err != nil {
		fmt.Println("error at decoding")
		return nil, err
	}
	return checkups, nil
}

func (c *MongoCheckUpStore) GetCheckByID(ctx context.Context, id string) (*types.CheckUp, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var checkup *types.CheckUp
	if err := c.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&checkup); err != nil {
		return nil, err
	}
	return checkup, nil

}

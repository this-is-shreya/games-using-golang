package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Ctx context.Context
var Client *mongo.Client
var Cancel context.CancelFunc
var IsConnected bool

func Close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func Connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error, bool) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	Ctx = context.Background()
	_, Cancel = context.WithTimeout(context.Background(),
		30*time.Second)

	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(uri))
	Client = client
	IsConnected = true
	fmt.Println("- ", IsConnected)
	return Client, Ctx, Cancel, err, IsConnected

}

// This is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func Ping(client *mongo.Client, ctx context.Context) error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

// func Query(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {

// 	// select database and collection.
// 	collection := client.Database(dataBase).Collection(col)

// 	// collection has an method Find,
// 	// that returns a mongo.cursor
// 	// based on query and field.
// 	result, err = collection.Find(ctx, query,
// 		options.Find().SetProjection(field))
// 	return
// }

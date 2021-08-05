package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const databaseName string = "testdb"

const mongoConnectionStringTest string = "mongodb://127.0.0.1:27017/?readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"
const countLenTest int = 10000

var authenticationCollectionTest *mongo.Collection

func VerificationFirst() {
	////connect MongoDB
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionStringTest))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	authenticationCollectionTest = client.Database(databaseName).Collection("authentication")
	for i := 0; i < countLenTest; i++ {
		wg.Add(1)
		go func(_wg *sync.WaitGroup, _i int) {

			updateResult, err := authenticationCollectionTest.UpdateOne(context.Background(), bson.M{"email": fmt.Sprintf("%d@.", _i)}, bson.D{primitive.E{Key: "$set", Value: bson.M{"is_activated": true}}})
			if err != nil {
				fmt.Printf("Can't set activated with %v and error: %v", updateResult, err)
			}
			defer _wg.Done()
		}(&wg, i)
	}
	wg.Wait()

	defer cancel()
	client.Disconnect(context.TODO())
}

func DropCollection() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionStringTest))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	authenticationCollectionTest = client.Database(databaseName).Collection("authentication")

	var wg sync.WaitGroup
	for i := 0; i < countLenTest; i++ {
		wg.Add(1)
		go func(_wg *sync.WaitGroup, _i int) {
			_, err := authenticationCollectionTest.DeleteOne(ctx, bson.M{"email": fmt.Sprintf("%d@.", _i)})
			if err != nil {
				log.Fatal(err)
			}
			defer _wg.Done()
		}(&wg, i)

	}
	wg.Wait()
	defer cancel()
	client.Disconnect(context.TODO())
}

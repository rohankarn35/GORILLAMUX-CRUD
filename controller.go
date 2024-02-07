package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}

var userCollection = db().Database("goTest").Collection("users")

func createProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person user
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		log.Fatal(err)
	}
	insertResult, err := userCollection.InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(insertResult.InsertedID)
}

func getUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body user
	errs := json.NewDecoder(r.Body).Decode(&body)
	if errs != nil {
		log.Fatal(errs)
	}
	var result primitive.M
	err := userCollection.FindOne(context.TODO(), bson.D{{Key: "name", Value: body.Name}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(result)

}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type updateBody struct {
		Name string `json:"name"`
		City string `json:"city"`
	}

	var body updateBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.D{{Key: "name", Value: body.Name}}
	after := options.After
	reurnOption := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: body.Name}}}}
	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, &reurnOption)
	var result primitive.M
	_ = updateResult.Decode(&result)
	json.NewEncoder(w).Encode(result)

}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"]
	_id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		log.Fatal(err)
	}
	opts := options.Delete().SetCollation(&options.Collation{})
	res, err := userCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(res.DeletedCount)

}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M
	cur, err := userCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {

		fmt.Println(err)

	}
	for cur.Next(context.TODO()) {

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}
	cur.Close(context.TODO())
	json.NewEncoder(w).Encode(results)
}

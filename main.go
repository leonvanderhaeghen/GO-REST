package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/leonvanderhaeghen/go-rest/helper"
	"github.com/leonvanderhaeghen/go-rest/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Connection mongoDB with helper class
var collection = helper.ConnectDB()

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Book array
	var products []models.Product

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var product models.Product
		// & character returns the memory address of the following variable.
		err := cur.Decode(&product) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		books = append(products, product)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(products) // encode similar to serialize process.
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var product models.Product
	// we get params with mux.
	var params = mux.Vars(r)


	id, _ := (params["id"])

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&product)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var product models.Product

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&product)

	// insert our book model.
	result, err := collection.InsertOne(context.TODO(), product)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := params["id"]

	var product models.Product

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&product)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"name", product.Name},
			{"price", product.Price},
			{"tags", product.Tags},
			{"barcode", product.Barcode},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&product)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	book.ID = id

	json.NewEncoder(w).Encode(book)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := params["id"]

	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

// var client *mongo.Client

func main() {
	//Init Router
	r := mux.NewRouter()

	r.HandleFunc("/api/products", getProducts).Methods("GET")
	r.HandleFunc("/api/product/{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/product", createProduct).Methods("POST")
	r.HandleFunc("/api/product/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/api/product/{id}", DeleteProduct).Methods("DELETE")

	config := helper.GetConfiguration()
	log.Fatal(http.ListenAndServe(config.Port, r))

}
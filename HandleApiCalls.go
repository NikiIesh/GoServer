package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	//"reflect"
	"strings"
	"time"

	// "regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name string `json:"name,omitempty"`
	Roll string `json:"empid,omitempty"`
}

func DBConnect() (context.Context, *mongo.Collection, *mongo.Client) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	ClientOptions, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = ClientOptions.Connect(ctx)
	if err != nil {
		fmt.Println("Mondodb Connection Error")
		os.Exit(1)
	}
	col := ClientOptions.Database("Sample_DB").Collection("Json_Dumps")
	return ctx, col, ClientOptions
}

func storage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var context, collection, ClientOptions = DBConnect()
	defer ClientOptions.Disconnect(context)
	switch r.Method {
	case "GET":
		IdFromURL := r.URL.Path
		DataID := strings.TrimLeft(IdFromURL, "/*")
		// col := client.Database("Sample_DB").Collection("Json_Dumps")
		objID, _ := primitive.ObjectIDFromHex(DataID)
		value := collection.FindOne(context, bson.M{"_id": objID})
		var bson_obj bson.M
		if err2 := value.Decode(&bson_obj); err2 != nil {
			fmt.Println(err2)
		}
		fmt.Println(bson_obj)

	case "POST":
		var person Person
		json.NewDecoder(r.Body).Decode(&person)
		DBConnect()
		defer r.Body.Close()
		document, err := collection.InsertOne(context, person)
		if err != nil {
			fmt.Println("Insert Error")
			os.Exit(1)
		} else {
			fmt.Println("Data Inserted Successfully")
			newID := document.InsertedID
			fmt.Println("ID: ", newID)
		}
	}
}

func main() {
	fmt.Println("Server is running")
	http.HandleFunc("/^[a-zA-Z0-9]{24}", storage)
	http.HandleFunc("/", storage)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	log "github.com/sirupsen/logrus"
	"net/http"
	"context"
	"time"
	//"github.com/gorilla/mux"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	id          string
	name        string
	description string
	picture     string
	priceUsd map[string]interface{}
	categories []string
}

type Products struct {
	Products []Product
}

var (
	productMutex *sync.Mutex
	product      Product
	myproduct   [] interface{}
)

var result map[string]interface{}
//var result Products

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	productMutex = &sync.Mutex{}

}

func main() {
	result := readFile()
	fmt.Println(result["products"])
	log.Info("Starting the http server")
	db_connection()
	// addProduct()
	handleRequests()
	
}
//this is my routers. 
func handleRequests(){
	http.HandleFunc("/", homePage)
	http.HandleFunc("/products", returnAllProducts)
	http.HandleFunc("/insertproduct", insertProduct)
	http.HandleFunc("/getproduct", getProduct)
	fmt.Println("Listening to http://localhost:10000/")
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllProducts(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(result["products"])
}


func insertProduct(w http.ResponseWriter, r *http.Request){
	//fmt.Println(r.Body)
	//reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("this is the body of the request", r.URL.Query())
	client := db_connection()
	var product Product
	json.Unmarshal(reqBody, &product)
	// send back the product it self. 
	collection := client.Database("products").Collection("salesproducts")
	// first try to find the product and then try to append it. 
	res, insertErr := collection.InsertOne(context.Background(), bson.D{product})
	
	json.NewEncoder(w).Encode(product)
}

func getProduct(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	id = params["id"]
	client := db_connection()
	collection := client.Database("products").Collection("salesproducts")
	//search through the collection for that particular id

	fmt.Println("this is the body of the request", r.URL.Query())
}

func db_connection()(*mongo.Client){
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb0.example.com:27017"))
	if err != nil{
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil{
		log.Fatal(err)
	}
	return client

}


// func addProduct(){
// 	client := db_connection()
// 	collection := client.Database("products").Collection("salesproducts")
// 	doc := result["products"]
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	res, insertErr := collection.InsertOne(ctx, doc)
// 	if insertErr != nil{
// 		log.Fatal(insertErr)
// 	}
// 	fmt.Println(res)
// 	fmt.Println(collection)

// }

func readFile() (map[string]interface{}) {
	if productMutex == nil {
		log.Fatalf("The productMutex is not initialized!")
		os.Exit(0)
	}
	productMutex.Lock()
	defer productMutex.Unlock()
	jsonFile, err := os.Open("products.json")
	defer jsonFile.Close()
    if err != nil {
        log.Fatalf("The jsonFile is not initialized")
    }
    log.Info("Successfully Opened product.json")
    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal([]byte(byteValue), &result)
    return result
}


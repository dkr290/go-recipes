package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	db = os.Getenv("MONGO_DATABASE")
	if db == "" {
		log.Fatal("You must set your 'MONGODB_DATABASE' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	createdb := os.Getenv("INITIAL_CREATEDB")
	if createdb == "" {
		log.Fatal("You must set your 'INITIAL_CREATEDB' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	if cl, err := mongo.Connect(ctx, options.Client().ApplyURI(uri)); err != nil {
		if err != nil {
			log.Fatal(err)
		}
	} else {
		MClient = cl
	}

	if err := MClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Mongo DB server")

	if createdb == "true" {
		initialDbCreate()

	}

}

func initialDbCreate() {

	type RecipeString struct {
		//swagger:ignore
		ID           string    `json:"id"`
		Name         string    `json:"name"`
		Tags         []string  `json:"tags"`
		Ingredients  []string  `json:"ingredients"`
		Instructions []string  `json:"instructions"`
		PublishedAt  time.Time `json:"publishedAt"`
	}

	recipes := make([]RecipeString, 0)
	file, _ := os.ReadFile("recipes.json")

	_ = json.Unmarshal([]byte(file), &recipes)
	var listOfRecipes []interface{}
	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}

	collection := MClient.Database(db).Collection("recipes")

	inseartManyResult, err := collection.InsertMany(context.Background(), listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("inserted recipes:", len(inseartManyResult.InsertedIDs))
}

func RedisConnect() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	redis_host := os.Getenv("REDIS_HOST")
	if redis_host == "" {
		log.Fatal("You must set your 'REDIS_HOST' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redis_host + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	status := redisClient.Ping(context.Background())
	fmt.Println(status)
}

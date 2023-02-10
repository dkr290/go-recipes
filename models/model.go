package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// swagger:parameters recipes newRecipe
type Recipe struct {
	//swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

var redisClient *redis.Client
var MClient *mongo.Client
var db string
var ctx = context.Background()

func GetAll() []Recipe {

	rr := make([]Recipe, 0)
	return rr
}

func GetRecipe() Recipe {

	var r1 Recipe
	return r1
}

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

func ListRecipes() ([]Recipe, error) {

	collection := MClient.Database(db).Collection("recipes")
	val, err := redisClient.Get(ctx, "recipes").Result()

	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		cur, err := collection.Find(ctx, bson.M{})
		if err != nil {
			return nil, err
		}
		defer cur.Close(ctx)
		recipes := make([]Recipe, 0)
		for cur.Next(ctx) {
			var recipe Recipe
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}
		data, _ := json.Marshal(recipes)
		redisClient.Set(ctx, "recipes", string(data), 0)
		return recipes, nil

	} else if err != nil {
		return nil, err
	} else {
		log.Printf("Request to Redis")
		recipes := make([]Recipe, 0)
		json.Unmarshal([]byte(val), &recipes)
		return recipes, nil
	}

}

func NewRecipe(c *gin.Context) (Recipe, error) {
	collection := MClient.Database(db).Collection("recipes")

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return recipe, err
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := collection.InsertOne(ctx, recipe)
	if err != nil {
		return recipe, err
	}
	log.Println("Remove data from Redis")
	redisClient.Del(ctx, "recipes")
	return recipe, nil
}

func UpdateRecipe(c *gin.Context) (Recipe, error) {

	id := c.Param("id")
	var recipe Recipe
	collection := MClient.Database(db).Collection("recipes")

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return recipe, err

	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}

	_, err := collection.UpdateOne(ctx, filter, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"tags", recipe.Tags},
		{"ingredients", recipe.Ingredients},
	}}})
	if err != nil {
		return recipe, err
	}
	log.Println("Remove data from Redis")
	redisClient.Del(ctx, "recipes")
	return recipe, nil
}

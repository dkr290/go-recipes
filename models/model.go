package models

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	recipes := make([]Recipe, 0)
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

func ListRecipes() ([]Recipe, error) {

	collection := MClient.Database(db).Collection("recipes")

	curr, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer curr.Close(ctx)
	recipes := make([]Recipe, 0)
	for curr.Next(ctx) {
		var recipe Recipe
		curr.Decode(&recipe)
		recipes = append(recipes, recipe)

	}

	return recipes, nil

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
	return recipe, nil
}

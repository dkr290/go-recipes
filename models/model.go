package models

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// func GetAll() []Recipe {

// 	rr := make([]Recipe, 0)
// 	return rr
// }

// func GetRecipe() Recipe {

// 	var r1 Recipe
// 	return r1
// }

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

func DeleteRecipe(c *gin.Context) (*mongo.DeleteResult, error) {
	id := c.Param("id")
	collection := MClient.Database(db).Collection("recipes")

	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	log.Println("Remove data from Redis")
	redisClient.Del(ctx, "recipes")
	return result, nil
}

func FindSingleRecipe(c *gin.Context) (Recipe, error) {
	id := c.Param("id")
	collection := MClient.Database(db).Collection("recipes")
	var recipe Recipe

	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}

	err := collection.FindOne(context.TODO(), filter).Decode(&recipe)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return recipe, err
		}

	}
	return recipe, nil

}

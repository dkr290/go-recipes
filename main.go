// Recipes API
//
// This is a sample recipes API.
//
//	   Schemes: http
//  Host: localhost:8080
//	   BasePath: /
//	   Version: 1.0.0
//	   Contact: ggg
//
//	   Consumes:
//	   - application/json
//
//	   Produces:
//	   - application/json
// swagger:meta

package main

import (
	"github.com/dkr290/go-recipes/handlers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var err error
var client *mongo.Client

func init() {

	// recipes := models.GetAll()
	// file, _ := os.ReadFile("recipes.json")
	// _ = json.Unmarshal([]byte(file), &recipes)
	// var listOfRecipes []interface{}
	// for _, recipe := range recipes {
	// 	listOfRecipes = append(listOfRecipes, recipe)
	// }
	// collection := client.Database(db).Collection("recipes")
	// log.Println(collection)
	// inseartManyResult, err := collection.InsertMany(context.Background(), listOfRecipes)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("inserted recipes:", len(inseartManyResult.InsertedIDs))

}

func main() {

	handler := handlers.NewHandlers()

	router := gin.Default()
	router.POST("/recipes", handler.NewRecipeHandler)
	router.GET("/recipes", handler.ListRecipesHandler)
	router.PUT("/recipes/:id", handler.UpdateRecipeHandler)
	router.GET("/recipes/search", handler.SearchRecipesHandler)
	router.DELETE("/recipes/:id", handler.DeleteRecipehandler)

	router.Run("127.0.0.1:8080")

}

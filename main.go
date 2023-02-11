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
	"github.com/dkr290/go-recipes/models"
	"github.com/gin-gonic/gin"
)

func main() {

	handler := handlers.NewHandlers()
	models.Connect()
	models.RedisConnect()
	router := gin.Default()
	router.POST("/recipes", handler.NewRecipeHandler)
	router.GET("/recipes", handler.ListRecipesHandler)
	router.PUT("/recipes/:id", handler.UpdateRecipeHandler)
	router.GET("/recipes/search", handler.SearchRecipesHandler)
	router.DELETE("/recipes/:id", handler.DeleteRecipehandler)
	router.GET("/recipes/:id", handler.SearchSingleRecipehandler)

	router.Run("127.0.0.1:8080")

}

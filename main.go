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
	authorized := router.Group("/")
	authHandler := handlers.NewAuthHandler()

	router.GET("/recipes", handler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.GET("/recipes/search", handler.SearchRecipesHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", handler.NewRecipeHandler)
		authorized.GET("/recipes/:id", handler.SearchSingleRecipehandler)
		authorized.PUT("/recipes/:id", handler.UpdateRecipeHandler)

		authorized.DELETE("/recipes/:id", handler.DeleteRecipehandler)
	}

	router.Run("127.0.0.1:8080")

}

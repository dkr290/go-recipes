package main

import (
	"github.com/dkr290/go-recipes/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	handler := handlers.NewHandlers()

	router := gin.Default()
	router.POST("/recipes", handler.NewRecipeHandler)
	router.GET("/recipes", handler.ListRecipesHandler)
	router.PUT("/recipes/:id", handler.UpdateRecipeHandler)
	router.Run("127.0.0.1:8080")

}

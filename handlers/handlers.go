package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dkr290/go-recipes/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Handler struct{}

var recipes = models.GetAll()

func init() {

	file, err := os.ReadFile("recipes.json")
	if err != nil {
		log.Fatalf("Error read recipes")

	}
	if err = json.Unmarshal([]byte(file), &recipes); err != nil {
		log.Fatalf("Error unmarshalling the recipes.json")
	}

}

func NewHandlers() *Handler {

	return &Handler{}
}

func (h *Handler) NewRecipeHandler(c *gin.Context) {

	recipe := models.GetRecipe()

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)

}

func (h *Handler) ListRecipesHandler(c *gin.Context) {

	c.JSON(http.StatusOK, recipes)
}

func (h *Handler) UpdateRecipeHandler(c *gin.Context) {

	id := c.Param("id")
	recipe := models.GetRecipe()

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes[index] = recipe
	recipes[index].ID = id
	c.JSON(http.StatusOK, recipe)

}

func (h *Handler) DeleteRecipehandler(c *gin.Context) {

	id := c.Param("id")
	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}
	// this is how to delete excluding only the index the item
	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe with id " + id + " has been deleted",
	})

}

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

var recipes = make([]models.Recipe, 0)

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

	var recipe models.Recipe

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

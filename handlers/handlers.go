package handlers

import (
	"net/http"
	"time"

	"github.com/dkr290/go-recipes/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Handler struct{}

var recipes = make([]models.Recipe, 0)

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

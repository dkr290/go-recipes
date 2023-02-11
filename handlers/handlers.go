package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dkr290/go-recipes/models"
	"github.com/gin-gonic/gin"
)

type Handler struct{}

var recipes = models.GetAll()

// func init() {

// 	file, err := os.ReadFile("recipes.json")
// 	if err != nil {
// 		log.Fatalf("Error read recipes")

// 	}
// 	if err = json.Unmarshal([]byte(file), &recipes); err != nil {
// 		log.Fatalf("Error unmarshalling the recipes.json")
// 	}

// }

func NewHandlers() *Handler {

	return &Handler{}
}

// swagger:operation POST /recipes recipes newRecipe
// Create a new recipe
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input

func (h *Handler) NewRecipeHandler(c *gin.Context) {

	recipe, err := models.NewRecipe(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while inserting a new recipe",
		})
		return
	}
	c.JSON(http.StatusOK, recipe)

}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
func (h *Handler) ListRecipesHandler(c *gin.Context) {

	recipes, err := models.ListRecipes()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, recipes)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	'400':
//	    description: Invalid input
//	'404':
//	    description: Invalid recipe ID
func (h *Handler) UpdateRecipeHandler(c *gin.Context) {

	recipe, err := models.UpdateRecipe(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been updated: " + recipe.Name,
	})
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete an existing recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
//	'404':
//	    description: Invalid recipe ID
func (h *Handler) DeleteRecipehandler(c *gin.Context) {

	result, err := models.DeleteRecipe(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	f := fmt.Sprint(result.DeletedCount)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted: " + f,
	})

}

// swagger:operation GET /recipes/search recipes findRecipe
// Search recipes based on tags
// ---
// produces:
// - application/json
// parameters:
//   - name: tag
//     in: query
//     description: recipe tag
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
func (h *Handler) SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]models.Recipe, 0)
	recipes, err := models.ListRecipes()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "There was search werror occured " + err.Error(),
		})
		return
	}
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)

}

func (h *Handler) SearchSingleRecipehandler(c *gin.Context) {

	recipe, err := models.FindSingleRecipe(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "There was no document with this id found " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

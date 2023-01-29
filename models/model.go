package models

import (
	"time"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

func GetAll() []Recipe {

	rr := make([]Recipe, 0)
	return rr
}

func GetRecipe() Recipe {

	var r1 Recipe
	return r1
}

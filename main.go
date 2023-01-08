package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/vijeyash1/server/models"
)

var recipes []models.Recipe

func init() {
	recipes = []models.Recipe{}
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal(file, &recipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHAndler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)
	router.Run(":8080")
}

func NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	err := c.BindJSON(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": recipe})
}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": recipes})
}

func UpdateRecipeHAndler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	err := c.BindJSON(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	index := -1

	for i, recipe := range recipes {
		if recipe.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "recipe not found"})
		return
	}
	recipes[index] = recipe
	c.JSON(http.StatusOK, gin.H{"status": "successfully updated the recipe", "data": recipe})
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i, recipe := range recipes {
		if recipe.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "recipe not found"})
		return
	}
	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{"status": "successfully deleted the recipe"})
}

func SearchRecipeHandler(c *gin.Context) {
	query := c.Query("tags")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "query is required"})
		return
	}
	var results []models.Recipe
	found := false
	for i, recipe := range recipes {

		for _, tag := range recipe.Tags {
			if strings.Contains(tag, query) {
				found = true
				break
			}
		}
		if found {
			results = append(results, recipes[i])
		}
		found = false
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": results})
}



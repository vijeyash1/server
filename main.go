package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

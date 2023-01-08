package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/vijeyash1/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewRecipesHandler(collection *mongo.Collection, ctx context.Context, redisClient *redis.Client) *RecipesHandler {
	return &RecipesHandler{collection, ctx, redisClient}
}

func (handler *RecipesHandler) GetOneRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	var recipe models.Recipe
	err := handler.collection.FindOne(handler.ctx, bson.M{"_id": objectId}).Decode(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": recipe})
}

func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	err := c.BindJSON(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err = handler.collection.InsertOne(context.Background(), recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error inserting the value", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "_id": recipe.ID})
}

func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	curr, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	defer curr.Close(context.Background())
	var recipes []models.Recipe
	for curr.Next(context.Background()) {
		var recipe models.Recipe
		err := curr.Decode(&recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "unable to Unmarshal check the model", "message": err.Error()})
			return
		}
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": recipes})
}

func (handler *RecipesHandler) UpdateRecipeHAndler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	err := c.BindJSON(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)

	_, err = handler.collection.UpdateOne(handler.ctx, bson.M{"_id": objectId}, bson.M{"$set": bson.M{"name": recipe.Name, "ingredients": recipe.Ingredients, "tags": recipe.Tags, "instructions": recipe.Instructions}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "successfully updated the recipe", "data": recipe})
}

func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectId, _ := primitive.ObjectIDFromHex(id)

	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{"_id": objectId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "successfully deleted the recipe"})
}

func (handler *RecipesHandler) SearchRecipeHandler(c *gin.Context) {
	query := c.Query("tags")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "query is required"})
		return
	}
	curr, err := handler.collection.Find(handler.ctx, bson.M{"tags": bson.M{"$in": []string{query}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	defer curr.Close(handler.ctx)
	var results []models.Recipe
	for curr.Next(handler.ctx) {
		var recipe models.Recipe
		err := curr.Decode(&recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "unable to Unmarshal check the model", "message": err.Error()})
			return
		}
		results = append(results, recipe)
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": results})
}

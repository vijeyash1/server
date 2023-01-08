package routes

import (
	"context"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/vijeyash1/server/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		claims := &handlers.Claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid || tkn == nil {
				c.JSON(401, gin.H{"status": "error", "message": "Invalid token"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func InitRoutes() {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisresult, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(redisresult)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:password@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	handler := handlers.NewRecipesHandler(client.Database("recipes").Collection("recipes"), ctx, redisClient)

	router := gin.Default()
	router.GET("/recipes", handler.ListRecipesHandler)
	authHandler := &handlers.AuthHandler{}
	router.POST("/signin", authHandler.LoginHandler)
	router.POST("/refresh", authHandler.RefreshTokenHandler)
	Authorized := router.Group("/")
	Authorized.Use(AuthMiddleware())
	{
		Authorized.POST("/recipes", handler.NewRecipeHandler)
		Authorized.PUT("/recipes/:id", handler.UpdateRecipeHAndler)
		Authorized.DELETE("/recipes/:id", handler.DeleteRecipeHandler)
		Authorized.GET("/recipes/search", handler.SearchRecipeHandler)
		Authorized.GET("/recipes/:id", handler.GetOneRecipeHandler)
	}
	router.Run(":8080")
}

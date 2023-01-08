package routes

import (
	"context"
	"log"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/vijeyash1/server/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token := session.Get("sessionToken")
		if token == nil {
			c.JSON(401, gin.H{"status": "error", "message": "Please Signup if you are a new user and Login if you are an existing user"})
			c.Abort()
			return
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

	handler := handlers.NewRecipesHandler(client.Database("recipes").Collection("recipes"), client.Database("users").Collection("recipes"), ctx, redisClient)

	router := gin.Default()
	store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.GET("/recipes", handler.ListRecipesHandler)
	router.POST("/signin", handler.LoginHandler)
	router.POST("/refresh", handler.RefreshTokenHandler)
	router.POST("/signup", handler.SignupHandler)
	router.GET("/logout", handler.LogoutHandler)
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

package handlers

import (
	"crypto/sha256"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/vijeyash1/server/models"
	"go.mongodb.org/mongo-driver/bson"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}
func (handler *RecipesHandler) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"status": "Logged out successfully"})
}

func (handler *RecipesHandler) SignupHandler(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	h := sha256.New()

	user.Password = string(h.Sum([]byte(user.Password)))

	_, err = handler.usercollection.InsertOne(handler.ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error inserting the value", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (handler *RecipesHandler) RefreshTokenHandler(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Token is not expired yet"})
		return
	}
	expirationTime := time.Now().Add(time.Minute * 5)
	claims.ExpiresAt = expirationTime.Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": JWTOutput{Token: tokenString, Expires: expirationTime}})
}

func (handler *RecipesHandler) LoginHandler(c *gin.Context) {
	user := models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	h := sha256.New()
	curr := handler.usercollection.FindOne(handler.ctx, bson.M{
		"username": user.Username,
		"password": string(h.Sum([]byte(user.Password))),
	})
	if curr.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid username or password"})
		return
	}

	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("sessionToken", sessionToken)	
	session.Set("username", user.Username)
	session.Save()
	c.String(http.StatusOK, "Logged in successfully")
}

package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var adminUsername string
var adminPassword string
var SECRET_KEY string

func init() {
	godotenv.Load()

	adminUsername = os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		panic("ADMIN_USERNAME is not set")
	}

	adminPassword = os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		panic("ADMIN_PASSWORD is not set")
	}

	SECRET_KEY = os.Getenv("SECRET_KEY")
	if SECRET_KEY == "" {
		panic("SECRET_KEY is not set")
	}
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != adminUsername || password != adminPassword {
		c.JSON(401, gin.H{"message": "authentication failed"})
		return
	}

	// トークン生成
	token, err := generateToken(username)
	if err != nil {
		c.JSON(500, gin.H{"message": "server error"})
		return
	}

	session := sessions.Default(c)
	session.Set("token", token)
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func generateToken(userID string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour).Unix(),
	})
	return token.SignedString([]byte(SECRET_KEY))
}

func middleware(c *gin.Context) {
	if !Authorized(c) {
		c.Redirect(http.StatusFound, "/need_to_login")
		c.Abort()
		return
	}

	c.Next()
}

func Authorized(c *gin.Context) bool {
	session := sessions.Default(c)
	token_ := session.Get("token")
	if token_ == nil {
		return false
	}
	tokenStr, ok := token_.(string)
	if !ok {
		return false
	}
	if tokenStr == "" {
		return false
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return false
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return false
	}

	return true
}

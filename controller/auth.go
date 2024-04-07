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

func generateToken(userID string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour).Unix(),
	})
	return token.SignedString([]byte(SECRET_KEY))
}

// func middleware(c *gin.Context) {
// 	session := sessions.Default(c)
// 	token_ := session.Get("token").(string)
// 	if token_ == "" {
// 		c.JSON(403, gin.H{"message": "need to login"})
// 		c.Abort()
// 		return
// 	}
// 	token, err := jwt.Parse(token_, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(SECRET_KEY), nil
// 	})
// 	if err != nil {
// 		c.JSON(403, gin.H{"message": "need to login"})
// 		c.Abort()
// 		return
// 	}
// 	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
// 		c.JSON(403, gin.H{"message": "need to login"})
// 		c.Abort()
// 	}

// 	c.Next()
// }

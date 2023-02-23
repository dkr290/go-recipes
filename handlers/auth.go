package handlers

import (
	"net/http"
	"time"

	"github.com/dkr290/go-recipes/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/xid"
)

type AuthHandler struct{}
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (handler *AuthHandler) SignInHandler(c *gin.Context) {

	models.GetUserPass(c)

	// expirationTime := time.Now().Add(10 * time.Minute)

	// claims := &Claims{
	// 	Username: user.Username,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: expirationTime.Unix(),
	// 	},
	// }

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// tokenString, err := token.SignedString([]byte(
	// 	os.Getenv("JWT_SECRET")))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError,
	// 		gin.H{
	// 			"error": err.Error(),
	// 		})
	// 	return
	// }

	// jwtOutput := JWTOutput{
	// 	Token:   tokenString,
	// 	Expires: expirationTime,
	// }

	// c.JSON(http.StatusOK, jwtOutput)

	c.JSON(http.StatusOK, gin.H{
		"message": "user signed in",
	})
}

func (handler *AuthHandler) SignOutHandler(c *gin.Context) {

	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "signed out",
	})
}

func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)
		sessionToken := session.Get("token")

		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Not Logged",
			})

			c.Abort()
		}

		c.Next()

		// tokenValue := c.GetHeader("Authorization")
		// claims := &Claims{}
		// tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(t *jwt.Token) (interface{}, error) {
		// 	return []byte(os.Getenv("JWT_SECRET")), nil
		// })

		// if err != nil {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// }

		// if tkn == nil || !tkn.Valid {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// }

	}
}

func (handler *AuthHandler) RefreshHandler(c *gin.Context) {

	session := sessions.Default(c)
	sessionToken := session.Get("token")
	sessionUser := session.Get("username")

	if sessionToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid session cookie",
		})
	}
	sessionToken = xid.New().String()
	session.Set("username", sessionUser.(string))
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, gin.H{

		"message": "New session issued",
	})

}

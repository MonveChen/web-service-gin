/*
 * @Author: Monve
 * @Date: 2023-07-24 11:45:43
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-25 09:14:39
 * @FilePath: /web-service-gin/controllers/user/user.go
 */
package user

import (
	"fmt"
	"net/http"
	"time"
	"web-service-gin/utils/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Replace these with your own secret key and expiration time.
const (
	secretKey      = "your_secret_key"
	tokenExpiresIn = time.Hour * 2
)

type LoginRequest struct {
	Key     string `json:"key" example:"dodo"`
	Secrert string `json:"secrert" example:"721c6ff80a6d3e4ad4ffa52a04c60085"`
}

// @BasePath /api/v1

// PingExample godoc
// @Summary 登陆
// @Schemes
// @Description 登陆获取token
// @Tags User
// @Accept json
// @Produce json
// @Param user body LoginRequest true "Login information"
// @Router /user/login [post]
func LoginHandler(c *gin.Context) {
	var jsonData LoginRequest
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Replace this with your own login logic.
	// Here, we assume that the user is authenticated and a valid username is received.
	key := jsonData.Key
	secret := jsonData.Secrert

	role := "normal"

	//从数据库校验密码是否正确，此处简略
	if key == "dodo" && secret == "721c6ff80a6d3e4ad4ffa52a04c60085" {
		role = "owner"
	} else if key == "test" && secret == "test_secret" {
		role = "normal"
	} else {
		c.JSON(500, gin.H{"error": "key or secret not right"})
		return
	}

	// Create a new token with the username as a custom claim.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.CustomClaims{
		Username: key,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpiresIn).Unix(),
		},
	})

	// Sign the token with the secret key to get the complete encoded token as a string.
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

// @BasePath /api/v1

// @Security BearerAuth
// @Summary 登出
// @Schemes
// @Description 登出账号（将token放入黑名单）
// @Tags User
// @Accept json
// @Produce json
// @Router /user/logout [post]
func LogoutHandler(c *gin.Context) {

	token := c.GetHeader("Authorization")
	expiresAt, exists := c.Get("expiresAt")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get expiresAt"})
		return
	}

	expirationTime2 := time.Now()
	fmt.Print(expirationTime2)
	expirationTime := time.Unix(expiresAt.(int64), 0)
	if err := auth.AddBlack(token, expirationTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to blacklist token"})
		return
	}

	c.String(http.StatusOK, "Logged out successfully")
}

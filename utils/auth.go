package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"goDemo/models"
	"strings"
	"time"
)

type Auth struct {
	SecretKey string
}

// GenerateToken 生成JWT token
func (s *Auth) GenerateToken(user *models.UserModel) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(s.SecretKey))
}

// ParseToken 获取并解析JWT token，获取其中UserID
func (s *Auth) ParseToken(ctx *gin.Context) (uint, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return 0, errors.New("missing Authorization header")
	}

	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		return 0, errors.New("invalid Authorization header format")
	}

	tokenString := splitToken[1]
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.SecretKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return 0, errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID in token")
	}

	return uint(userID), nil
}

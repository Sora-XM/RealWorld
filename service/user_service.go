package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"goDemo/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserService struct {
	DB        *gorm.DB
	SecretKey string
}

// CreateUser 注册用户
func (s *UserService) CreateUser(user *models.UserModel) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	fmt.Printf("存储的哈希密码: %s\n", user.Password)
	return s.DB.Create(user).Error
}

// VerifyUser 验证用户信息
func (s *UserService) VerifyUser(email string, password string) (*models.UserModel, error) {
	var user models.UserModel
	err := s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, err
	}
	// 去除哈希密码两端的空白字符
	trimmedPassword := strings.TrimSpace(user.Password)
	fmt.Printf("从数据库读取的哈希密码（去除空白后）: %s\n", trimmedPassword)
	fmt.Printf("待验证的明文密码: %s\n", password)

	// 使用 bcrypt 验证密码
	//err = bcrypt.CompareHashAndPassword([]byte(trimmedPassword), []byte(password))
	//if err != nil {
	//	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
	//		return nil, nil // 密码错误
	//	}
	//	return nil, err // 其他 bcrypt 错误
	//}

	return &user, nil
}

// GenerateToken 生成JWT token
func (s *UserService) GenerateToken(user *models.UserModel) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(s.SecretKey))
}

// GetUserByToken 根据token获取用户信息
func (s *UserService) GetUserByToken(token string) (*models.UserModel, error) {
	// 解析 JWT 令牌
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	// 检查令牌是否有效
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	// 从令牌中获取用户 ID
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	// 根据用户 ID 从数据库中查询用户信息
	var user models.UserModel
	err = s.DB.First(&user, uint(userID)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

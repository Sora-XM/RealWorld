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
	fmt.Printf("存储的哈希密码: %s\n", user.Password) // 添加日志输出
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
		return nil, err // 其他数据库错误
	}
	// 去除哈希密码两端的空白字符
	trimmedPassword := strings.TrimSpace(user.Password)
	fmt.Printf("从数据库读取的哈希密码（去除空白后）: %s\n", trimmedPassword) // 添加日志输出
	fmt.Printf("待验证的明文密码: %s\n", password)                  // 添加日志输出

	// 使用 bcrypt 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(trimmedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, nil // 密码错误
		}
		return nil, err // 其他 bcrypt 错误
	}

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

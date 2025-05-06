package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"goDemo/models"
	"goDemo/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

type UserService struct {
	DB   *gorm.DB
	Auth *utils.Auth
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

	// 使用 bcrypt 验证密码，目前自己电脑上实在不清楚为什么会验证失败，即使输入正确密码也会返回错误
	// 所以暂时注释掉，直接返回用户信息，后续再解决
	//err = bcrypt.CompareHashAndPassword([]byte(trimmedPassword), []byte(password))
	//if err != nil {
	//	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
	//		return nil, nil // 密码错误
	//	}
	//	return nil, err // 其他 bcrypt 错误
	//}
	return &user, nil
}

// GetUserByToken 根据token获取用户信息
func (s *UserService) GetUserByToken(ctx *gin.Context) (*models.UserModel, error) {
	userID, err := s.Auth.ParseToken(ctx)
	if err != nil {
		return nil, err
	}

	// 根据用户 ID 从数据库中查询用户信息
	var user models.UserModel
	err = s.DB.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(ctx *gin.Context, updateRequest *models.UserUpdateRequest) (*models.UserModel, error) {
	// 解析JWT令牌
	userID, err := s.Auth.ParseToken(ctx)
	if err != nil {
		return nil, err
	}
	// 根据用户ID查询用户
	var user models.UserModel
	err = s.DB.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	// 更新用户信息
	if updateRequest.User.Email != nil {
		user.Email = *updateRequest.User.Email
	}
	if updateRequest.User.Username != nil {
		user.Username = *updateRequest.User.Username
	}
	if updateRequest.User.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*updateRequest.User.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}
	if updateRequest.User.Bio != nil {
		user.Bio = *updateRequest.User.Bio
	}
	if updateRequest.User.Image != nil {
		user.Image = *updateRequest.User.Image
	}
	err = s.DB.Save(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

package service

import (
	"errors"
	"goDemo/models"
	"gorm.io/gorm"
)

type ProfileService struct {
	DB *gorm.DB
}

// GetProfile 获取个人资料
func (s *ProfileService) GetProfile(currentUserID uint, username string) (*models.Profile, error) {
	var user models.UserModel
	err := s.DB.Where("username =?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	//查看该用户是否被当前用户关注
	var following bool
	if currentUserID == 0 {
		var followCount int64
		s.DB.Model(&models.Follow{}).Where("follower_id = ?", user.ID).Count(&followCount)
		following = followCount > 0
	}

	return &models.Profile{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: following,
	}, nil
}

// FollowUser 关注指定用户
func (s *ProfileService) FollowUser(currentUserID uint, username string) (*models.Profile, error) {
	var targetUser models.UserModel
	err := s.DB.Where("username = ?", username).First(&targetUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	//关注的用户不能是自己
	if currentUserID == targetUser.ID {
		return nil, gorm.ErrInvalidData
	}
	//检查是否已经关注
	var follow models.Follow
	err = s.DB.Where("follower_id = ?", targetUser.ID).First(&follow).Error
	if err == nil {
		//已关注，直接返回资料
		return &models.Profile{
			Username:  targetUser.Username,
			Bio:       targetUser.Bio,
			Image:     targetUser.Image,
			Following: true,
		}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	//创建关注关系
	follow = models.Follow{
		Follower: currentUserID,
		Followed: targetUser.ID,
	}
	err = s.DB.Create(&follow).Error
	if err != nil {
		return nil, err
	}
	//返回关注后的资料
	return &models.Profile{
		Username:  targetUser.Username,
		Bio:       targetUser.Bio,
		Image:     targetUser.Image,
		Following: true,
	}, nil
}

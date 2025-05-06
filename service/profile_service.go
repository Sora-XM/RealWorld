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
	if currentUserID != 0 {
		var follow models.Follow
		err = s.DB.Where("follower = ? AND followed = ?", currentUserID, user.ID).First(&follow).Error
		following = err == nil
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
	err = s.DB.Where("follower = ? AND followed = ?", currentUserID, targetUser.ID).First(&follow).Error
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

func (s *ProfileService) UnfollowUser(currentUserID uint, username string) (*models.Profile, error) {
	var targetUser models.UserModel
	err := s.DB.Where("username = ?", username).First(&targetUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	//是否已经关注
	var follow models.Follow
	err = s.DB.Where("follower = ? AND followed = ?", currentUserID, targetUser.ID).First(&follow).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Profile{
				Username:  targetUser.Username,
				Bio:       targetUser.Bio,
				Image:     targetUser.Image,
				Following: false,
			}, nil
		}
		return nil, err
	}
	//删除关注关系并返回资料
	err = s.DB.Delete(&follow).Error
	if err != nil {
		return nil, err
	}
	return &models.Profile{
		Username:  targetUser.Username,
		Bio:       targetUser.Bio,
		Image:     targetUser.Image,
		Following: false,
	}, nil
}

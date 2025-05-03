package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goDemo/service"
	"goDemo/utils"
	"gorm.io/gorm"
	"net/http"
)

type ProfileController struct {
	ProfileService *service.ProfileService
	UserService    *service.UserService
	Auth           *utils.Auth
}

// GetProfile godoc
// @Summary 获取用户资料
// @Description 根据用户名获取用户资料
// @Tags profiles
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   username path string true "用户名"
// @Success 200 {object} models.Profile "获取成功，返回用户资料"
// @Failure 401 {object} map[string]interface{} "未授权，缺少或无效的 token"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/profiles/{username} [get]
func (c *ProfileController) GetProfile(ctx *gin.Context) {
	username := ctx.Param("username")
	currentUserID, err := c.Auth.ParseToken(ctx)
	if err != nil {
		if err.Error() == "missing Authorization header" {
			// 没有 token 也可以继续获取资料，只是 currentUserID 为 0
			currentUserID = 0
			err = nil
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
			return
		}
	}
	// 调用服务层获取用户资料
	profile, err := c.ProfileService.GetProfile(currentUserID, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"body": []string{"用户不存在"}}})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{"服务器内部错误"}}})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"profile": profile})
}

// FollowUser godoc
// @Summary 关注用户
// @Description 关注指定的用户
// @Tags profiles
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   username path string true "用户名"
// @Success 200 {object} models.Profile "关注成功，返回用户资料"
// @Failure 401 {object} map[string]interface{} "未授权，缺少或无效的 token"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/profiles/{username}/follow [post]
func (c *ProfileController) FollowUser(ctx *gin.Context) {
	// 从路径参数获取要关注的用户名
	username := ctx.Param("username")
	currentUserID, err := c.Auth.ParseToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	if currentUserID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"未授权，缺少或无效的 token"}}})
		return
	}
	// 调用服务层关注用户
	profile, err := c.ProfileService.FollowUser(currentUserID, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"body": []string{"用户不存在"}}})
		} else if errors.Is(err, gorm.ErrInvalidData) {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"body": []string{"不能自己关注自己"}}})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		}
		return
	}
	// 返回关注后的用户资料
	ctx.JSON(http.StatusOK, gin.H{"profile": profile})
}

// UnfollowUser godoc
// @Summary 取消关注用户
// @Description 取消关注指定的用户
// @Tags profiles
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   username path string true "用户名"
// @Success 200 {object} models.Profile "取消关注成功，返回用户资料"
// @Failure 401 {object} map[string]interface{} "未授权，缺少或无效的 token"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/profiles/{username}/follow [delete]
func (c *ProfileController) UnfollowUser(ctx *gin.Context) {
	username := ctx.Param("username")
	currentUserID, err := c.Auth.ParseToken(ctx)
	//检查token
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	if currentUserID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"未授权，缺少或无效的 token"}}})
		return
	}
	//调用服务层取消关注用户
	profile, err := c.ProfileService.UnfollowUser(currentUserID, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"body": []string{"用户不存在"}}})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"profile": profile})
}

package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goDemo/service"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type ProfileController struct {
	ProfileService *service.ProfileService
	UserService    *service.UserService
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
	var currentUserID uint
	//获取认证信息，如果有token则验证token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) == 2 && strings.ToLower(splitToken[0]) == "bearer" {
			token := splitToken[1]
			user, err := c.UserService.GetUserByToken(token)
			if err != nil {
				currentUserID = user.ID
			}
		}
	}
	//调用服务层获取用户资料
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
	// 从请求头获取 Authorization 字段
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"缺少 Authorization 头"}}})
		return
	}
	// 提取 token
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"无效的 Authorization 头格式"}}})
		return
	}
	token := splitToken[1]
	// 通过 Token 获取当前用户信息
	user, err := c.UserService.GetUserByToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"无效的 token"}}})
		return
	}
	// 调用服务层关注用户
	profile, err := c.ProfileService.FollowUser(user.ID, username)
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

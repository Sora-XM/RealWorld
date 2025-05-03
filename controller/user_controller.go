package controller

import (
	"github.com/gin-gonic/gin"
	"goDemo/models"
	"goDemo/service"
	"goDemo/utils"
	"net/http"
	"strings"
)

type UserController struct {
	UserService *service.UserService
	Auth        *utils.Auth
}

// RegisterUser godoc
// @Summary 注册用户
// @Description 接收用户信息，将用户信息存储到数据库中完成注册
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user body models.UserModel true "用户注册信息"
// @Success 201 {object} models.UserModel "注册成功，返回创建的用户信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/users/register [post]
func (c *UserController) RegisterUser(ctx *gin.Context) {

	var registerRequest models.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := registerRequest.User
	if err := c.UserService.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// LoginUser godoc
// @Summary 用户登录
// @Description 接收用户登录信息，验证用户信息并返回token
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user body models.UserModel true "用户登录信息"
// @Success 200 {object} map[string]string "登录成功，返回token"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "用户名或密码错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/users/login [post]
func (c *UserController) LoginUser(ctx *gin.Context) {
	var loginRequest models.UserLoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}

	user, err := c.UserService.VerifyUser(loginRequest.User.Email, loginRequest.User.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"用户名或密码错误"}}})
		return
	}

	token, err := c.Auth.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	response := models.UserResponse{
		User: struct {
			Email    string `json:"email"`
			Token    string `json:"token"`
			Username string `json:"username"`
			Bio      string `json:"bio"`
			Image    string `json:"image"`
		}{
			Email:    user.Email,
			Token:    token,
			Username: user.Username,
			Bio:      user.Bio,
			Image:    user.Image,
		},
	}
	ctx.JSON(http.StatusOK, response)
}

// GetCurrentUser godoc
// @Summary 获取当前用户信息
// @Description 根据token获取当前用户信息
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} models.UserModel "获取成功，返回当前用户信息"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/user [get]
func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	user, err := c.UserService.GetUserByToken(ctx)
	if err != nil {
		if err.Error() == "missing Authorization header" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"缺少 Authorization 请求头"}}})
		} else if err.Error() == "invalid Authorization header format" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"无效的 Authorization 请求头形式"}}})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{"无效的Token"}}})
		}
		return
	}
	// 生成新的 token
	token, err := c.Auth.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	// 构建响应
	response := models.UserResponse{
		User: struct {
			Email    string `json:"email"`
			Token    string `json:"token"`
			Username string `json:"username"`
			Bio      string `json:"bio"`
			Image    string `json:"image"`
		}{
			Email:    user.Email,
			Token:    token,
			Username: user.Username,
			Bio:      user.Bio,
			Image:    user.Image,
		},
	}
	ctx.JSON(http.StatusOK, response)
}

// UpdateUser godoc
// @Summary 更新用户信息
// @Description 根据token更新当前用户信息
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   user body models.UserModel true "用户更新信息"
// @Success 200 {object} models.UserModel "更新成功，返回更新后的用户信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/user [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// 解析请求体
	var updateRequest models.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	// 调用服务层更新用户信息
	updatedUser, err := c.UserService.UpdateUser(ctx, &updateRequest)
	if err != nil {
		if err.Error() == "missing Authorization header" || err.Error() == "invalid Authorization header format" || strings.Contains(err.Error(), "invalid token") || strings.Contains(err.Error(), "user not found") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		}
		return
	}
	// 生成新的 token
	newToken, err := c.Auth.GenerateToken(updatedUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
		return
	}
	// 构建响应
	response := models.UserResponse{
		User: struct {
			Email    string `json:"email"`
			Token    string `json:"token"`
			Username string `json:"username"`
			Bio      string `json:"bio"`
			Image    string `json:"image"`
		}{
			Email:    updatedUser.Email,
			Token:    newToken,
			Username: updatedUser.Username,
			Bio:      updatedUser.Bio,
			Image:    updatedUser.Image,
		},
	}

	// 返回响应
	ctx.JSON(http.StatusOK, response)
}

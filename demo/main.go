package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
)

func sqlDriver() {
	dsh := "user:password@tcp(127.0.0.1:3306)/sql_demo"
	db, err := sql.Open("mysql", dsh)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func main() {
	////创建一个服务，根据url映射，发送消息或数据
	//ginServer := gin.Default()
	////加载静态资源
	//ginServer.LoadHTMLGlob("templates/*")
	////加载资源页面
	//ginServer.Static("/static", "./static")
	//
	//ginServer.GET("/user", func(context *gin.Context) {
	//	//context.JSON(200, gin.H{"message": "hello world"})
	//	context.HTML(http.StatusOK, "index.html", gin.H{
	//		"massage": "Go后台数据Sora",
	//	})
	//})
	//ginServer.POST("/hello", func(context *gin.Context) {
	//	context.JSON(200, gin.H{"message": "hello post"})
	//})
	//
	////user/info?name=Sora&age=18
	//ginServer.POST("/user/info", func(context *gin.Context) {
	//	//获取url参数
	//	name := context.Query("name")
	//	age := context.Query("age")
	//	context.JSON(200, gin.H{
	//		"name": name,
	//		"age":  age,
	//	})
	//})
	//ginServer.GET("/user/info/:name/:age", func(context *gin.Context) {
	//	//获取url参数
	//	name := context.Param("name")
	//	age := context.Param("age")
	//	context.JSON(200, gin.H{
	//		"name": name,
	//		"age":  age,
	//	})
	//})
	//
	////前端给后端传Json
	//ginServer.POST("/user/json", func(context *gin.Context) {
	//	data, _ := context.GetRawData()
	//	var m map[string]interface{}
	//	//包装为json数据
	//	_ = json.Unmarshal(data, &m)
	//	//获取json数据
	//	context.JSON(http.StatusOK, m)
	//})
	//
	////前端给后端传form表单
	//ginServer.POST("/user/form", func(context *gin.Context) {
	//	//获取form表单数据
	//	username := context.PostForm("username")
	//	password := context.PostForm("password")
	//	context.JSON(http.StatusOK, gin.H{
	//		"message":  "ok",
	//		"username": username,
	//		"password": password,
	//	})
	//})
	//
	////路由
	//ginServer.GET("/test", func(context *gin.Context) {
	//	//重定向 301
	//	context.Redirect(301, "http://hatsune.com.cn/")
	//})
	//
	//// 404
	//ginServer.NoRoute(func(context *gin.Context) {
	//	context.HTML(http.StatusNotFound, "404.html", gin.H{})
	//})
	//
	////使用路由组
	//userGroup := ginServer.Group("/user")
	//{
	//	userGroup.GET("/add")
	//	userGroup.GET("/login")
	//	userGroup.GET("/passage")
	//}
	//
	//ginServer.Run(":8082")
	// 推荐长度为 32 字节或更长
	keyLength := 32
	secretKey, err := generateSecretKey(keyLength)
	if err != nil {
		fmt.Println("生成 SecretKey 失败:", err)
		return
	}
	fmt.Println("生成的 SecretKey:", secretKey)
}
func generateSecretKey(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

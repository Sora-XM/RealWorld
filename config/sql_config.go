package config

import (
	_ "github.com/go-sql-driver/mysql"
	"goDemo/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	dsh := "root:740627119zxl@tcp(127.0.0.1:3306)/realworld_sql?charset=utf8&parseTime=true"
	db, err := gorm.Open(mysql.Open(dsh), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.UserModel{}, &models.Profile{}, &models.Follow{}, &models.Article{}, &models.Comment{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

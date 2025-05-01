package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 输入原始密码
	var originalPassword string
	fmt.Print("请输入原始密码: ")
	fmt.Scanln(&originalPassword)

	// 生成哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(originalPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("生成哈希失败:", err)
		return
	}
	fmt.Printf("生成的哈希密码: %s\n", string(hashedPassword))

	// 输入待验证的密码
	var verifyPassword string
	fmt.Print("请输入待验证的密码: ")
	fmt.Scanln(&verifyPassword)

	// 验证密码
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(verifyPassword))
	if err != nil {
		fmt.Println("密码验证失败:", err)
	} else {
		fmt.Println("密码验证成功")
	}
}

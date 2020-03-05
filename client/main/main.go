package main

import (
	"fmt"
	"go_code/chatroot/client/process"
	"os"
)

//定义两个变量，一个表示用户id, 一个表示用户密码
var userId int
var userPwd string
var userName string

func main()  {
	//接收用户的选择
	var key int
	//判断是否继续显示菜单
	//var loop = true
	for true {
		fmt.Println("-------------欢迎登陆多人聊天系统-------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择（1-3）: ")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码:")
			fmt.Scanf("%s\n", &userPwd)
			//完成登陆
			//1. 创建一个UserProcess实例
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id: ")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码: ")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户名字(nickname): ")
			fmt.Scanf("%s\n", &userName)
			//2. 调用UserProcess，完成注册请求
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)
		case 3:
			fmt.Println("退出系统")
			//loop = false
			os.Exit(0)
		default:
			fmt.Println("输入有误，请重新输入")
		}
	}

	//根据用户输入，显示新的提示
	/*
	if key == 1 {
		//说明用户要登陆
		fmt.Println("请输入用户id:")
		fmt.Scanf("%d\n", &userId)
		fmt.Println("请输入用户的密码:")
		fmt.Scanf("%s\n", &userPwd)

		//登陆函数，写到另外一个文件 login.go
		//这里需要重新调用
		//login(userId, userPwd)

			//if err != nil {
			//	fmt.Println("登陆失败")
			//} else {
			//	fmt.Println("登陆成功")
			//}

	} else if key == 2 {
		fmt.Println("进行用户注册")
	}
*/
}
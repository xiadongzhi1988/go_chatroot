package main

import (
	"fmt"
	"go_code/chatroot/server/model"
	"net"
	"time"
)


//处理和客户端的通讯
func process(conn net.Conn)  {
	defer conn.Close()
	//调用总控，创建一个总控
	processor := &Processor{
		conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器端通讯协程错误 err=", err)
		return
	}
}

func init()  {
	//当服务器启动时，就去初始化redis连接池
	initPool("127.0.0.1:6379", 16, 0, 300 * time.Second)
	initUserDao()
}

//编写一个函数，完成对UserDao的初始化任务
func initUserDao()  {
	//pool本身就是一个全局变量
	//初始化顺序问题
	model.MyUserDao = model.NewUserDao(pool)
}

func main()  {
	//提示信息
	fmt.Println("服务器[新的结构]在8889端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}
	//一旦监听成功，就等待客户端连接服务器
	for  {
		fmt.Println("等待客户端来连接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		//一旦连接成功，则启动一个协程和客户端保持通讯
		go process(conn)
	}
}
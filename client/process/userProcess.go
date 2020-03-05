package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/chatroot/client/utils"
	"go_code/chatroot/common/message"
	"net"
	"os"
)

type UserProcess struct {

}

func (this *UserProcess) Register(userId int,
	userPwd string, userName string) (err error) {
	//1. 连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	defer conn.Close()
	//2. 准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType
	// 3. 创建一个LoginMes 结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	//4. 将registerMes 序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal  err=", err)
		return
	}

	//5.把data赋给mes.data字段
	mes.Data = string(data)

	//6. 将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal  err=", err)
		return
	}

	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	//发送data给服务器
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
	}

	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}

	//将mes的data部分反序列化成 RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功，重新登录一下")
		os.Exit(0)
	} else  {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}

//给关联一个用户登陆的方法
//写一个函数，完成登陆
func (this *UserProcess) Login(userId int, userPwd string) (err error) {
	// 定协议
	//fmt.Printf("userId=%d userPwd=%s\n", userId, userPwd)
	//return nil
	//1. 连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	defer conn.Close()
	//2. 准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	// 3. 创建一个LoginMes 结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	//4. 将loginMes 序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal  err=", err)
		return
	}
	//5.把data赋给mes.data字段
	mes.Data = string(data)
	//6. 将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal  err=", err)
		return
	}
	//7. 这时data就是要发送的消息
	//7.1 先把data长度发送给服务器
	//先获取到data的长度-> 转成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	fmt.Printf("客户端发送消息长度=%d 内容=%s\n", len(data), string(data))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	//休眠20
	//time.Sleep(20 * time.Second)
	//fmt.Println("休眠了20..")
	//这里需要处理处理服务器端返回的消息
	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}
	//将mes的data部分反序列化成 LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		//fmt.Println("登陆成功")
		//显示当前在线用户列表，变量loginResMes.UsersId
		fmt.Println("当前在线用户列表如下:")
		for _, v := range loginResMes.UsersId {
			//不显示自己在线
			if v == userId {
				continue
			}
			fmt.Println("用户id:\t", v)
			//完成客户端 onlineUsers 完成初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")

		//客户端启动一个协程
		//该协程保持和服务器的通讯，如果服务器有数据推送给客户端
		//则接收并显示在客户端的终端
		go serverProcessMes(conn)

		//1. 显示登陆成功的菜单[循环]
		for  {
			ShowMenu()
		}
	} else  {
		fmt.Println(loginResMes.Error)
	}
	return
}
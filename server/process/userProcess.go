package process2

import (
	"encoding/json"
	"fmt"
	"go_code/chatroot/common/message"
	"go_code/chatroot/server/model"
	"go_code/chatroot/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加一耳光字段，表示该Conn是哪个用户
	UserId int
}

//通知所有在线的用户的方法
//userId 通知其它的在线用户，我上线
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历 onlineUsers， 然后一个一个发送 NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		//过滤掉自己
		if id == userId {
			continue
		}
		//开始通知[单独写一个方法]
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int) {
	//组装 NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline
	//将 notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//将序列化后的 notifyUserStatusMes 赋值给mes.Data
	mes.Data = string(data)
	//对mes再次序列化,准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//发送，创建Transfer实例
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
}

func (this UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//1. 先从mes中取出mes.Data,并直接反序列化成registerMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	//1.先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//到redis数据库去完成注册
	//1. 使用model.MyUserDao 到redis去验证
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误..."
		}
	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}

	//4. 将data赋值给 resMes
	resMes.Data = string(data)
	//5. 对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	//6. 发送data，将其封装到writePkg()函数
	//因为使用分层模式(mvc)， 先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

//编写一个函数serverProcessLogin函数，专门处理登陆请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//1. 先从mes中取出mes.Data,并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//1.先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//2.在声明LoginResMes, 并完成赋值
	var loginResMes message.LoginResMes
	//到redis数据库去完成验证
	//1. 使用model.MyUserDao 到redis去验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误.."
		}
	} else {
		loginResMes.Code = 200
		//这里用户登陆成功，把该登陆成功的用户放入userMgr中
		//将登陆成功的用户的userId赋给this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)
		//通知其它的在线的用户，我上线了
		this.NotifyOthersOnlineUser(loginMes.UserId)
		//将当前在线用户的id放入 loginResMes.UsersId
		//遍历userMgr。onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user, "登陆成功")
	}

	/*
	//如果用户id=100， 密码=123456， 认为合法
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
		//合法
		loginResMes.Code = 200
	} else {
		//不合法
		loginResMes.Code = 500 //500状态码，表示用户不存在
		loginResMes.Error = "用户不存在,请注册再使用..."
	}
	 */

	//3. 将loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	//4. 将data赋值给 resMes
	resMes.Data = string(data)
	//5. 对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	//6. 发送data，将其封装到writePkg()函数
	//因为使用分层模式(mvc)， 先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
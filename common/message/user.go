package message

//定义一个用户的结构体
type User struct {
	//为了能序列化和反序列化成功，必须保证
	//用户信息的json字符串的key 和结构体的字段对应的tag名字一致！
	UserId int `json:"userId"`
	UserPwd string `json:"userPwd"`
	UserName string `json:"userName"`
	UserStatus int `json:"userStatus"` //用户状态
	Sex string `json:"sex"` //性别
}
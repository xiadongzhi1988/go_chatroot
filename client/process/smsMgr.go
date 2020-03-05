package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroot/common/message"
)

func outputGroupMes(mes *message.Message)  { //这个地方mes一定是SmsMes
	//显示
	//1. 反序列化mes.Data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err.Error())
		return
	}
	//显示信息
	info := fmt.Sprintf("用户id:\t%d 对大家说:\t%s",
		smsMes.UserId, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}
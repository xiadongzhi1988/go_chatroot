package model

import (
	"go_code/chatroot/common/message"
	"net"
)
//在客户端很多地方会使用到curUser，作为一个全局
type CurUser struct {
	Conn net.Conn
	message.User
}
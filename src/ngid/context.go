/**
 * @Author: mjzheng
 * @Description:
 * @File:  context.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 下午6:32
 */

package ngid

import (
	"github.com/mjproto/simple_msg"
	"net"
)

type SimpleMsgContext struct {
	Conn    net.Conn
	HeadReq *simple_msg.HeadReq
	HeadRsp *simple_msg.HeadRsp
	RawData []byte
}

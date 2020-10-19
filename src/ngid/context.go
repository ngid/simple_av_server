/**
 * @Author: mjzheng
 * @Description:
 * @File:  context.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 下午6:32
 */

package ngid

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/mjproto/simple_msg"
	"google.golang.org/grpc"
	"net"
)

type SimpleMsgContext struct {
	Conn    net.Conn
	HeadReq *simple_msg.HeadReq
	HeadRsp *simple_msg.HeadRsp

	BodyReq proto.Message
	BodyRsp proto.Message
	RawData []byte

	Stream grpc.ServerStream
}

func GetSimpleContext(ctx context.Context) *SimpleMsgContext {
	msgContext := ctx.Value("ngid").(*SimpleMsgContext)
	return msgContext
}
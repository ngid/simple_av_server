/**
 * @Author: mjzheng
 * @Description:
 * @File:  trpc_server.go
 * @Version: 1.0.0
 * @Date: 2020/10/16 下午8:39
 */

package ngid

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/mjproto/simple_msg"
	"io"
	"log"
	"reflect"
)

// 业务实现方法的容器
type TrpcServer struct {
}

func (s *TrpcServer) Head(gs simple_msg.SimpleMsg_HeadServer) error {
	for {
		headReq, err := gs.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}

		headRsp := &simple_msg.HeadRsp{
			Cmd:    headReq.GetCmd(),
			Subcmd: headReq.GetSubcmd(),
			Seq:    headReq.GetSeq(),
		}

		ctx1 := context.Background()
		ctx, _ := context.WithCancel(ctx1)
		ngid := &SimpleMsgContext{
			Conn:    nil,
			HeadReq: headReq,
			HeadRsp: headRsp,
			stream:  gs,
		}
		connCtx := context.WithValue(ctx, "ngid", ngid)

		HandleTrpcMsg(connCtx)
	}
	return nil
}

func HandleTrpcMsg(ctx context.Context) {

	//fmt.Println(headReq)

	msgContext := GetSimpleContext(ctx)
	headReq := msgContext.HeadReq
	headRsp := msgContext.HeadRsp
	reqBodyType, rspBodyType, handler, err := GetRegisterFunc(headReq.Cmd, headReq.Subcmd)
	if err != nil {
		return
	}

	bodyReq, ok := reflect.New(reqBodyType.Elem()).Interface().(proto.Message)
	if !ok {
		return
	}

	bodyRsp, ok := reflect.New(rspBodyType.Elem()).Interface().(proto.Message)
	if !ok {
		return
	}

	if err := proto.Unmarshal(headReq.GetEx(), bodyReq); err != nil {
		return
	}

	msgContext.BodyReq = bodyReq
	msgContext.BodyRsp = bodyRsp

	headRsp.ErrCode, headRsp.ErrMsg = handler.HandleMsg(ctx)

	//fmt.Println(msgContext.BodyReq, msgContext.BodyRsp)
	log.Printf("headReq[%v] req[%v] headRsp[%v] rsp[%v]", msgContext.HeadReq, msgContext.BodyReq, msgContext.HeadRsp, msgContext.BodyRsp)

	msgContext.stream.SendMsg(bodyRsp)
}

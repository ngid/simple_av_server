/**
 * @Author: mjzheng
 * @Description:
 * @File:  parse.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 下午7:02
 */

package ngid

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mjproto/simple_msg"
	"log"
	"reflect"
)

const (
	STATUS_START_EX = 1
	STATUS_LENGTH   = 2
	STATUS_BODY     = 3
	STATUS_END_EX   = 4
	STATUS_COMPLETE = 5
)

func ParseMsg(ctx context.Context, buf []byte, total int) (remain []byte, remainLen int) {
	useLen := 0
	from := 0
	status := STATUS_START_EX
	needLen := 1
	for from+needLen <= total {
		switch status {
		case STATUS_START_EX:
			if buf[from] != 0x2 {
				fmt.Println("unexcept start error")
				break
			}
			from += needLen
			needLen = 2
			status = STATUS_LENGTH
		case STATUS_LENGTH:
			msgLen := int(binary.BigEndian.Uint16(buf[from : from+needLen]))
			from += needLen
			needLen = msgLen
			status = STATUS_BODY
		case STATUS_BODY:
			HandleMsg(ctx, buf[from:from+needLen])
			from += needLen
			needLen = 1
			status = STATUS_END_EX
		case STATUS_END_EX:
			if buf[from] != 0x3 {
				fmt.Println("unexcept end error")
				break
			}
			from += needLen
			needLen = 1
			status = STATUS_COMPLETE
		case STATUS_COMPLETE:
			useLen = from
			status = STATUS_START_EX
		}
	}

	if useLen < total {
		// move
		remainLen = total - useLen
		for i := 0; i < remainLen; i++ {
			buf[i] = buf[useLen+i]
		}
		//fmt.Println("reamin len", total, useLen, remainLen)
		return buf, remainLen
	} else {
		return buf, 0
	}
}

func HandleMsg(ctx context.Context, pData []byte) {
	headReq := &simple_msg.HeadReq{}
	err := proto.Unmarshal(pData, headReq)
	if err != nil {
		//panic(err)
		fmt.Println(err)
		return
	}

	//fmt.Println(headReq)
	headRsp := &simple_msg.HeadRsp{
		Cmd:    headReq.GetCmd(),
		Subcmd: headReq.GetSubcmd(),
		Seq:    headReq.GetSeq(),
	}

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

	msgContext := GetSimpleContext(ctx)
	msgContext.HeadReq = headReq
	msgContext.HeadRsp = headRsp
	msgContext.RawData = pData
	msgContext.BodyReq = bodyReq
	msgContext.BodyRsp = bodyRsp

	headRsp.ErrCode, headRsp.ErrMsg = handler.HandleMsg(ctx)

	//fmt.Println(msgContext.BodyReq, msgContext.BodyRsp)
	log.Printf("headReq[%v] req[%v] headRsp[%v] rsp[%v]", msgContext.HeadReq, msgContext.BodyReq, msgContext.HeadRsp, msgContext.BodyRsp)

	pRsp := ComposeMsg(headRsp)
	msgContext.Conn.Write(pRsp)
}

func ComposeMsg(msg proto.Message) (data []byte) {
	pData, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	buf.WriteByte(0x2)
	lenBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBuf, uint16(len(pData)))
	buf.Write(lenBuf)
	buf.Write(pData)
	buf.WriteByte(0x3)

	//fmt.Println(len(pData))

	data = buf.Bytes()
	return
}

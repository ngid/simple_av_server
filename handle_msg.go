/**
 * @Author: mjzheng
 * @Description:
 * @File:  handle_msg.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 上午11:28
 */

package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mjproto/simple_av"
	"github.com/mjproto/simple_msg"
	"net"
)

func HandleMsg(ctx context.Context, pData []byte) {
	msg := &simple_msg.HeadReq{}
	err := proto.Unmarshal(pData, msg)
	if err != nil {
		//panic(err)
		fmt.Println(err)
		return
	}

	fmt.Println(msg)
	headRsp := &simple_msg.HeadRsp{
		Cmd:    msg.GetCmd(),
		Subcmd: msg.GetSubcmd(),
		Seq:    msg.GetSeq(),
	}

	switch msg.Subcmd {
	case int32(simple_av.SUB_CMD_JoinRoom):
		req := &simple_av.JoinRoomReq{}
		proto.Unmarshal(msg.Ex, req)
		rsp := &simple_av.JoinRoomRsp{}
		HandleJoinRoom(ctx, req, rsp)
		headRsp.Ex, _ = proto.Marshal(rsp)
		//fmt.Println(req)
	case int32(simple_av.SUB_CMD_ExitRoom):
		req := &simple_av.ExitRoomReq{}
		proto.Unmarshal(msg.Ex, req)
		rsp := &simple_av.ExitRoomRsp{}
		HandleExitRoom(ctx, req, rsp)
		headRsp.Ex, _ = proto.Marshal(rsp)
		//fmt.Println(req)
	case int32(simple_av.SUB_CMD_Upload):
		req := &simple_av.UploadReq{}
		proto.Unmarshal(msg.Ex, req)
		rsp := &simple_av.UploadRsp{}
		headRsp.ErrCode, headRsp.ErrMsg = HandleUpload(ctx, req, rsp)
		headRsp.Ex, _ = proto.Marshal(rsp)
		//fmt.Println(req)
	case int32(simple_av.SUB_CMD_SendData):
		req := &simple_av.SendDataReq{}
		proto.Unmarshal(msg.Ex, req)
		rsp := &simple_av.SendDataRsp{}
		headRsp.ErrCode, headRsp.ErrMsg = HandleSendData(ctx, req, rsp)
		headRsp.Ex, _ = proto.Marshal(rsp)
	}

	pRsp, _ := proto.Marshal(headRsp)
	conn := ctx.Value("conn").(net.Conn)
	conn.Write(pRsp)
}

func HandleJoinRoom(ctx context.Context, req *simple_av.JoinRoomReq, rsp *simple_av.JoinRoomRsp) (errorCode int32, errorMsg string) {

	uid := req.GetUid()
	_, ok := UidToRoomId[uid]
	if ok {
		return 1, "already join room"
	}

	//
	roomId := req.GetRoomId()
	UidToRoomId[uid] = roomId //

	roomInfo, ok := RoomList[roomId]
	if !ok {
		roomInfo = RoomInfo{
			userList: make(map[int64]UserInfo),
		}
		RoomList[roomId] = roomInfo
	}
	conn := ctx.Value("conn").(net.Conn)
	roomInfo.AddUser(uid, conn)

	return 0, ""
}

func HandleExitRoom(ctx context.Context, req *simple_av.ExitRoomReq, rsp *simple_av.ExitRoomRsp) (errorCode int32, errorMsg string) {
	uid := req.GetUid()
	if _, ok := UidToRoomId[uid]; ok {
		delete(UidToRoomId, uid)
	}

	roomId := req.GetRoomId()
	if roomInfo, ok := RoomList[roomId]; ok {
		roomInfo.DeleteUser(uid)
	}
	return 0, ""
}

func HandleUpload(ctx context.Context, req *simple_av.UploadReq, rsp *simple_av.UploadRsp) (errorCode int32, errorMsg string) {
	uid := req.GetUid()
	roomId := req.GetRoomId()
	if roomInfo, ok := RoomList[roomId]; ok {
		roomInfo.UpdateUser(uid, true)
	}
	return 0, ""
}

func HandleSendData(ctx context.Context, req *simple_av.UploadReq, rsp *simple_av.UploadRsp) (errorCode int32, errorMsg string) {
	return 0, ""
}

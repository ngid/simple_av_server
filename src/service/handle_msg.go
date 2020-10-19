/**
 * @Author: mjzheng
 * @Description:
 * @File:  handle_msg.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 上午11:28
 */

package service

import (
	"context"
	"github.com/mjproto/simple_av"
	"github.com/ngid/simple_av_server/src/ngid"
)

func HandleJoinRoom(ctx context.Context) (errorCode int32, errorMsg string) {
	msgContext := ngid.GetSimpleContext(ctx)
	req := msgContext.BodyReq.(*simple_av.JoinRoomReq)
	//rsp := msgContext.BodyRsp.(*simple_av.JoinRoomRsp)

	roomId := req.GetRoomId()
	uid := req.GetUid()

	roomInfo := RManager.GetRoom(ctx, roomId)
	conn := msgContext.Conn
	gs := msgContext.Stream
	roomInfo.AddUser(uid, conn, gs)

	return 0, "success"
}

func HandleExitRoom(ctx context.Context) (errorCode int32, errorMsg string) {
	msgContext := ngid.GetSimpleContext(ctx)

	req := msgContext.BodyReq.(*simple_av.ExitRoomReq)

	roomId := req.GetRoomId()
	uid := req.GetUid()

	roomInfo := RManager.GetRoom(ctx, roomId)
	roomInfo.DeleteUser(uid)
	return 0, "success"
}

func HandleUpload(ctx context.Context) (errorCode int32, errorMsg string) {
	msgContext := ngid.GetSimpleContext(ctx)
	req := msgContext.BodyReq.(*simple_av.UploadReq)

	roomId := req.GetRoomId()
	uid := req.GetUid()
	roomInfo := RManager.GetRoom(ctx, roomId)
	roomInfo.UpdateUser(uid, true)
	return 0, "success"
}

func HandleSendData(ctx context.Context) (errorCode int32, errorMsg string) {
	msgContext := ngid.GetSimpleContext(ctx)
	req := msgContext.BodyReq.(*simple_av.SendDataReq)

	roomId := req.GetRoomId()
	uid := req.GetUid()
	roomInfo := RManager.GetRoom(ctx, roomId)
	//roomInfo.SendAll(uid, msgContext.RawData)
	roomInfo.SendAllUseTRPC(uid)
	return 0, "success"
}

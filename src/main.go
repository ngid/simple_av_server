/**
 * @Author: mjzheng
 * @Description:
 * @File:  main.go
 * @Version: 1.0.0
 * @Date: 2020/6/22 下午8:02
 */

package main

import (
	"github.com/mjproto/simple_av"
	"github.com/ngid/simple_av_server/src/ngid"
	"github.com/ngid/simple_av_server/src/service"
)

func init() {
	ngid.RegisterFunc(int32(simple_av.BIG_CMD_SIMPLE_AV), int32(simple_av.SUB_CMD_JoinRoom), &simple_av.JoinRoomReq{}, &simple_av.JoinRoomRsp{},
		ngid.HandlerFunc(service.HandleJoinRoom))

	ngid.RegisterFunc(int32(simple_av.BIG_CMD_SIMPLE_AV), int32(simple_av.SUB_CMD_ExitRoom), &simple_av.ExitRoomReq{}, &simple_av.ExitRoomRsp{},
		ngid.HandlerFunc(service.HandleExitRoom))

	ngid.RegisterFunc(int32(simple_av.BIG_CMD_SIMPLE_AV), int32(simple_av.SUB_CMD_Upload), &simple_av.UploadReq{}, &simple_av.UploadRsp{},
		ngid.HandlerFunc(service.HandleUpload))

	ngid.RegisterFunc(int32(simple_av.BIG_CMD_SIMPLE_AV), int32(simple_av.SUB_CMD_SendData), &simple_av.SendDataReq{}, &simple_av.SendDataRsp{},
		ngid.HandlerFunc(service.HandleSendData))

}

func main() {
	Listen()
}

/**
 * @Author: mjzheng
 * @Description:
 * @File:  main.go
 * @Version: 1.0.0
 * @Date: 2020/6/22 下午8:02
 */

package main

import (
	"fmt"
	"github.com/mjproto/simple_av"
	"github.com/ngid/simple_av_server/src/ngid"
	"github.com/ngid/simple_av_server/src/service"
	"reflect"
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

func TestReflect() {
	//req := simple_av.JoinRoomReq{
	//	RoomId: 100,
	//	Uid: 88,
	//}
	////joinType := reflect.TypeOf(req)
	//joinValue := reflect.ValueOf(req)
	//
	//fmt.Println(joinValue.Kind())
	//fmt.Println(joinValue.Type())
	//
	//fmt.Println("start get ")
	//
	//for i:=0; i<joinValue.NumField(); i++ {
	//	fmt.Println(joinValue.Field(i))
	//}

	joinType := reflect.TypeOf(&simple_av.JoinRoomReq{})

	fmt.Println(joinType, joinType.Elem())

	valueType := reflect.New(joinType.Elem())
	fmt.Println(valueType)

	req := valueType.Interface().(*simple_av.JoinRoomReq)
	fmt.Printf("%#v", req)

	//fmt.Println(joinType, joinValue)
}

func main() {
	//TestReflect()
	ngid.Listen()
}

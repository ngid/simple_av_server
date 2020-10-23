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
	"github.com/mjproto/simple_msg"
	_ "github.com/ngid/simple_av_server/src/log_files"
	"github.com/ngid/simple_av_server/src/ngid"
	"github.com/ngid/simple_av_server/src/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strings"
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o simple_av_server

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

func WordCount(s string) map[string]int {
	m := make(map[string]int)
	words := strings.Split(s, " ")
	for _, word := range  words {
		if v, ok := m[word]; ok {
			m[word] = v + 1
		} else {
			m[word] = 1
		}
	}
	return m
}


func main() {
	// ngid.Listen("localhost:50000")

	lis, err := net.Listen("tcp", ":50000") //监听所有网卡8028端口的TCP连接
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	s := grpc.NewServer() //创建gRPC服务

	/**注册接口服务
	 * 以定义proto时的service为单位注册，服务中可以有多个方法
	 * (proto编译时会为每个service生成Register***Server方法)
	 * 包.注册服务方法(gRpc服务实例，包含接口方法的结构体[指针])
	 */

	simple_msg.RegisterSimpleMsgServer(s, &ngid.TrpcServer{})
	/**如果有可以注册多个接口服务,结构体要实现对应的接口方法
	 * user.RegisterLoginServer(s, &server{})
	 * minMovie.RegisterFbiServer(s, &server{})
	 */
	// 在gRPC服务器上注册反射服务
	reflection.Register(s)
	// 将监听交给gRPC服务处理
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

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
	ngid.RegisterFunc(0x66f, 3, &simple_av.JoinRoomReq{}, &simple_av.JoinRoomRsp{}, ngid.HandlerFunc(service.HandleJoin))
}

func main() {
	Listen()
}

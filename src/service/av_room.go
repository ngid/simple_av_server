/**
 * @Author: mjzheng
 * @Description:
 * @File:  av_room.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 下午2:25
 */

package service

import (
	"context"
	"fmt"
	"github.com/ngid/simple_av_server/src/ngid"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type UserInfo struct {
	uid     int64
	conn    net.Conn
	bUpload bool
	stream  grpc.ServerStream
}

type UserInfoList map[int64]*UserInfo

type RoomInfo struct {
	roomId   int64
	userList UserInfoList // uid->UserInfo
	mutex    sync.Mutex
}

func (c *RoomInfo) AddUser(uid int64, conn net.Conn, gs grpc.ServerStream) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.userList[uid]; !ok {
		userInfo := &UserInfo{
			uid:     uid,
			conn:    conn,
			bUpload: false,
			stream:  gs,
		}
		c.userList[uid] = userInfo
		log.Println("add user", userInfo)
	}
}

func (c *RoomInfo) DeleteUser(uid int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.userList[uid]; ok {
		delete(c.userList, uid)
	}
	c.mutex.Unlock()
}

func (c *RoomInfo) UpdateUser(uid int64, bUpload bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if user, ok := c.userList[uid]; ok {
		fmt.Printf("%p\n", &user)
		user.bUpload = bUpload
	}
}

func (c *RoomInfo) SendAll(uid int64, b []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	data := ngid.ComposeMsgWithBytes(b)
	for _, user := range c.userList {
		if user.uid == uid {
			continue
		}
		user.conn.Write(data)
	}
}

func (c *RoomInfo) SendAllUseTRPC(ctx context.Context, uid int64) {
	msgContext := ngid.GetSimpleContext(ctx)
	req := msgContext.HeadReq
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, user := range c.userList {
		if user.uid == uid {
			continue
		}

		user.stream.SendMsg(req)
	}
}

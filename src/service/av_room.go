/**
 * @Author: mjzheng
 * @Description:
 * @File:  av_room.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 下午2:25
 */

package service

import (
	"fmt"
	"github.com/ngid/simple_av_server/src/ngid"
	"net"
	"sync"
)

type UserInfo struct {
	uid     int64
	conn    net.Conn
	bUpload bool
}

type UserInfoList map[int64]*UserInfo

type RoomInfo struct {
	roomId   int64
	userList UserInfoList // uid->UserInfo
	mutex    sync.Mutex
}

func (c *RoomInfo) AddUser(uid int64, conn net.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.userList[uid]; !ok {
		userInfo := &UserInfo{
			uid:     uid,
			conn:    conn,
			bUpload: false,
		}
		c.userList[uid] = userInfo
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

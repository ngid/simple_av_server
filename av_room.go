/**
 * @Author: mjzheng
 * @Description:
 * @File:  av_room.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 下午2:25
 */

package main

import (
	"net"
	"sync"
)

type UserInfo struct {
	uid     int64
	conn    net.Conn
	bUpload bool
}

type RoomInfo struct {
	userList map[int64]UserInfo // uid->UserInfo
	mutex    sync.Mutex
}

func (c *RoomInfo) AddUser(uid int64, conn net.Conn) {
	c.mutex.Lock()
	if _, ok := c.userList[uid]; !ok {
		userInfo := UserInfo{
			uid:     uid,
			conn:    conn,
			bUpload: false,
		}
		c.userList[uid] = userInfo
	}
	c.mutex.Unlock()
}

func (c *RoomInfo) DeleteUser(uid int64) {
	c.mutex.Lock()
	if _, ok := c.userList[uid]; ok {
		delete(c.userList, uid)
	}
	c.mutex.Unlock()
}

func (c *RoomInfo) UpdateUser(uid int64, bUpload bool) {
	c.mutex.Lock()
	if user, ok := c.userList[uid]; ok {
		user.bUpload = bUpload
	}
	c.mutex.Unlock()
}

func (c *RoomInfo) SendAll(b []byte) {
	c.mutex.Lock()
	for _, user := range c.userList {
		user.conn.Write(b)
	}
	c.mutex.Unlock()
}

var RoomList map[int64]RoomInfo // room id-> room info

var UidToRoomId map[int64]int64 // uid->room id

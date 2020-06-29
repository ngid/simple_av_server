/**
 * @Author: mjzheng
 * @Description:
 * @File:  av_room_manager.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 下午4:35
 */
package service

import (
	"context"
	"sync"
)

type RoomInfoList map[int64]*RoomInfo // 类型别名

type RoomManager struct {
	RoomList RoomInfoList // room id-> room info
	mutex    sync.Mutex
}

var RManager RoomManager = RoomManager{
	RoomList: make(RoomInfoList),
	mutex:    sync.Mutex{},
}

func (m *RoomManager) GetRoom(ctx context.Context, roomId int64) *RoomInfo {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	roomInfo, ok := m.RoomList[roomId]
	if !ok {
		roomInfo = &RoomInfo{
			roomId:   roomId,
			userList: make(UserInfoList),
			mutex:    sync.Mutex{},
		}
		m.RoomList[roomId] = roomInfo
	}
	return roomInfo
}

func (m *RoomManager) DeleteRoom(ctx context.Context, roomId int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.RoomList, roomId)
}

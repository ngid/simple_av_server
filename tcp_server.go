package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
)

func Listen() {
	listener, err := net.Listen("tcp", "localhost:50000")
	if err != nil {
		return
	}

	ctx1 := context.Background()
	ctx, _ := context.WithCancel(ctx1)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		connCtx := context.WithValue(ctx, "conn", conn)
		//clients.Add(conn)
		go doServerStuff(conn, connCtx)
	}
}

func doServerStuff(conn net.Conn, ctx context.Context) {
	buf := make([]byte, 1024)
	fmt.Println("len", len(buf))
	from := 0
	for {
		total, err := conn.Read(buf[from:])
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return //终止程序
		}

		buf, from = SpiltPackage(ctx, buf, total+from)
	}
}

const (
	STATUS_START_EX = 1
	STATUS_LENGTH   = 2
	STATUS_BODY     = 3
	STATUS_END_EX   = 4
	STATUS_COMPLETE = 5
)

func SpiltPackage(ctx context.Context, buf []byte, total int) (remain []byte, remainLen int) {
	useLen := 0
	from := 0
	status := STATUS_START_EX
	needLen := 1
	for from+needLen <= total {
		switch status {
		case STATUS_START_EX:
			if buf[from] != 0x2 {
				fmt.Println("unexcept start error")
				break
			}
			from += needLen
			needLen = 2
			status = STATUS_LENGTH
		case STATUS_LENGTH:
			msgLen := int(binary.BigEndian.Uint16(buf[from : from+needLen]))
			from += needLen
			needLen = msgLen
			status = STATUS_BODY
		case STATUS_BODY:
			HandleMsg(ctx, buf[from:from+needLen])
			from += needLen
			needLen = 1
			status = STATUS_END_EX
		case STATUS_END_EX:
			if buf[from] != 0x3 {
				fmt.Println("unexcept end error")
				break
			}
			from += needLen
			needLen = 1
			status = STATUS_COMPLETE
		case STATUS_COMPLETE:
			useLen = from
			status = STATUS_START_EX
		}
	}

	if useLen < total {
		// move
		remainLen = total - useLen
		for i := 0; i < remainLen; i++ {
			buf[i] = buf[useLen+i]
		}
		//fmt.Println("reamin len", total, useLen, remainLen)
		return buf, remainLen
	} else {
		return buf, 0
	}
}

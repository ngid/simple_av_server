package main

import (
	"context"
	"fmt"
	"github.com/ngid/simple_av_server/src/ngid"
	"github.com/ngid/simple_av_server/src/service"
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
		ngid := &ngid.SimpleMsgContext{
			Conn: conn,
		}
		connCtx := context.WithValue(ctx, "ngid", ngid)
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

		buf, from = service.ParseMsg(ctx, buf, total+from)
	}
}

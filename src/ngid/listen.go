package ngid

import (
	"context"
	"fmt"
	"net"
)

func Listen(addr string) {
	listener, err := net.Listen("tcp", addr)
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
		ngid := &SimpleMsgContext{
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

		buf, from = ParseMsg(ctx, buf, total+from)
	}
}

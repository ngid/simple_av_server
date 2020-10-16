package ngid

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var ConnectionList []net.Conn

func Listen(addr string) {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("failed to listen")
		return
	}

	go func() {
		c := <-stopChan
		log.Println("user stop listen", c)
		if err = listener.Close(); err != nil {

		}
		for _, conn := range ConnectionList {
			conn.Close()
		}
	}()

	ctx1 := context.Background()
	ctx, _ := context.WithCancel(ctx1)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("user stop listen, exit accept")
			break
		}
		ngid := &SimpleMsgContext{
			Conn: conn,
		}
		ConnectionList = append(ConnectionList, conn)
		connCtx := context.WithValue(ctx, "ngid", ngid)
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

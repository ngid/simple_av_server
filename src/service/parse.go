/**
 * @Author: mjzheng
 * @Description:
 * @File:  parse.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 下午7:02
 */

package service

import (
	"context"
	"encoding/binary"
	"fmt"
)

const (
	STATUS_START_EX = 1
	STATUS_LENGTH   = 2
	STATUS_BODY     = 3
	STATUS_END_EX   = 4
	STATUS_COMPLETE = 5
)

func ParseMsg(ctx context.Context, buf []byte, total int) (remain []byte, remainLen int) {
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

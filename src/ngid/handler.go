/**
 * @Author: mjzheng
 * @Description:
 * @File:  handler.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 下午7:44
 */

package ngid

import "context"

type Handler interface {
	HandleMsg(ctx context.Context) (errorCode int32, errorMsg string)
}

type HandlerFunc func(ctx context.Context) (errorCode int32, errorMsg string)

func (f HandlerFunc) HandleMsg(ctx context.Context) (errorCode int32, errorMsg string) {
	return f(ctx)
}

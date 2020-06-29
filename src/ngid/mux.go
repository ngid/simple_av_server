/**
 * @Author: mjzheng
 * @Description:
 * @File:  mux.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 下午7:44
 */

package ngid

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
	"runtime"
	"sync"
)

type iliveCmd struct {
	bigCmd uint32
	subCmd uint32
}

//String 用于输出
func (cmd iliveCmd) String() string {
	return fmt.Sprintf("0x%x:0x%x", cmd.bigCmd, cmd.subCmd)
}

type register struct {
	reqType reflect.Type
	rspType reflect.Type
	h       Handler
}

//String 用户输出
func (mux register) String() string {
	m := make(map[string]string, 4)
	m["ReqType"] = mux.reqType.String()
	m["RspType"] = mux.rspType.String()
	m["handler"] = runtime.FuncForPC(reflect.ValueOf(mux.h).Pointer()).Name()
	b, _ := json.MarshalIndent(m, "", "    ")
	return string(b)
}

type ServeMux struct {
	mu          sync.RWMutex
	m           map[string]muxEntry
	mapRegister map[iliveCmd]register
}

type muxEntry struct {
	h       Handler
	pattern string
}

var DefaultServeMux = &ServeMux{
	m:           make(map[string]muxEntry),
	mapRegister: make(map[iliveCmd]register),
}

func HandleFunc(pattern string, handler HandlerFunc) {
	DefaultServeMux.HandleFunc(pattern, handler)
}

func Handle(pattern string, handler Handler) {
	DefaultServeMux.Handle(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler HandlerFunc) {
	if handler == nil {
		panic("gifts: nil handler")
	}
	mux.Handle(pattern, handler)
}

func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if handler == nil {
		panic("gifts: nil handler")
	}

	if _, exist := mux.m[pattern]; exist {
		//panic("http: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	e := muxEntry{h: handler, pattern: pattern}
	mux.m[pattern] = e
}

func (mux *ServeMux) GetHandler(pattern string) Handler {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if v, ok := mux.m[pattern]; ok {
		return v.h
	}
	return nil
}

// RegisterFunc 注册ILIVE请求处理函数
// bigCmd 大命令字, subCmd 小命令字
// 最后注册覆盖先前注册
// handler为空，则panic
// req or rsp为空,则panic
func (mux *ServeMux) RegisterFunc(bigCmd uint32, subCmd uint32, req proto.Message, rsp proto.Message, handler HandlerFunc) {

	if bigCmd == 0 || subCmd == 0 {
		panic("ilive: invalid cmd")
	}

	if handler == nil {
		panic("ilive: nil handler")
	}

	if req == nil || rsp == nil {
		panic("ilive:req or rsp is nil")
	}

	cmd := iliveCmd{
		bigCmd: bigCmd,
		subCmd: subCmd,
	}

	mux.mu.Lock()
	defer mux.mu.Unlock()
	fmt.Println(reflect.TypeOf(req), reflect.TypeOf(rsp))
	mux.mapRegister[cmd] = register{reqType: reflect.TypeOf(req), rspType: reflect.TypeOf(rsp), h: HandlerFunc(handler)}
}

// RegisterFunc default mux handler
func RegisterFunc(bigCmd uint32, subCmd uint32, req proto.Message, rsp proto.Message, handler HandlerFunc) {
	DefaultServeMux.RegisterFunc(bigCmd, subCmd, req, rsp, handler)
}

// GetRegisterInfo 查找注册信息
func GetRegisterInfo(ctx context.Context, mux *ServeMux, bigCmd uint32, subCmd uint32) (reflect.Type, reflect.Type, Handler, error) {
	cmd := iliveCmd{
		bigCmd: bigCmd,
		subCmd: subCmd,
	}
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if reg, ok := mux.mapRegister[cmd]; ok {
		return reg.reqType, reg.rspType, reg.h, nil
	}
	return nil, nil, nil, nil

	//pattern := fmt.Sprintf("/0x%x/%d", bigCmd, subCmd)
	//entry, ok := mux.m[pattern]
	//if !ok {
	//	return nil, nil, nil, errors.New("Cmd Not Match")
	//}
	//
	//reqBodyName := entry.obj + "Req"
	//reqBodyType := proto.MessageType(reqBodyName)
	//if reqBodyType == nil {
	//	return nil, nil, nil, fmt.Errorf("proto MessageType Not Found: %s", reqBodyName)
	//}
	//
	//rspBodyName := entry.obj + "Rsp"
	//rspBodyType := proto.MessageType(rspBodyName)
	//if rspBodyType == nil {
	//	return nil, nil, nil, fmt.Errorf("proto MessageType Not Found: %s", rspBodyName)
	//}
	//return reqBodyType, rspBodyType, entry.h, nil
}

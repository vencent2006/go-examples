/**
 * @Author: vincent
 * @Description:
 * @File:  connection
 * @Version: 1.0.0
 * @Date: 2021/8/14 23:22
 */

package websocket

import (
	kim "go-examples/chat_demo_im"
	"net"

	"github.com/gobwas/ws"
)

// Frame
type Frame struct {
	raw ws.Frame
}

func (f *Frame) SetOpCode(code kim.OpCode) {
	f.raw.Header.OpCode = ws.OpCode(code)
}

func (f *Frame) GetOpCode() kim.OpCode {
	return kim.OpCode(f.raw.Header.OpCode)
}

func (f *Frame) SetPayload(payload []byte) {
	f.raw.Payload = payload
}

func (f *Frame) GetPayload() []byte {
	if f.raw.Header.Masked {
		// 解密
		ws.Cipher(f.raw.Payload, f.raw.Header.Mask, 0)
	}

	// todo: 为啥要吧Masked设置为false呢，感觉ws.Cipher就把f.raw.Payload进行了修改
	// 为防止二次执行cipher，所以第一次cipher完就设置为masked=false
	// 那并发时，不会出现2次cipher的问题吗
	f.raw.Header.Masked = false
	return f.raw.Payload
}

// websocket conn
type WsConn struct {
	net.Conn
}

func (c *WsConn) ReadFrame() (kim.Frame, error) {
	f, err := ws.ReadFrame(c.Conn)
	if err != nil {
		return nil, err
	}
	return &Frame{raw: f}, nil
}

func (c *WsConn) WriteFrame(code kim.OpCode, payload []byte) error {
	f := ws.NewFrame(ws.OpCode(code), true, payload)
	return ws.WriteFrame(c.Conn, f)
}

func (c *WsConn) Flush() error {
	// todo: 实现flush
	return nil
}

func NewWsConn(conn net.Conn) *WsConn {
	return &WsConn{conn}
}

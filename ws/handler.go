package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/websocket"
	"log"
)

var upgrader = websocket.HertzUpgrader{
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}
var wsMaxBufSize = 1024 * 1024 * 10
var wsPool = gopool.NewPool("ws", 30, gopool.NewConfig())

func Handler(_ context.Context, c *app.RequestContext) {
	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		handler := &wsConnHandler{conn: conn}
		handler.spin()
	})
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
}

type wsConnHandler struct {
	conn *websocket.Conn
	buf  bytes.Buffer
}

// ws 消息循环
func (h *wsConnHandler) spin() {
	for {
		_, message, err := h.conn.ReadMessage()
		if err != nil {
			hlog.Error("ws read:", err)
			break
		}
		hlog.Debugf("ws recv: %s", message)
		h.buf.Write(message)
		h.processRequests()
		if h.buf.Len() > wsMaxBufSize {
			hlog.Error("ws buf is full, just clear it")
			h.buf.Reset()
		}
	}

	_ = h.conn.Close()
}

// 处理 ws 消息
func (h *wsConnHandler) processRequests() {
	frame := ReadFrame(&h.buf)
	if frame == nil {
		return
	}

	var req Request
	err := json.Unmarshal(frame.Body, &req)
	if err != nil {
		hlog.Error("ws unmarshal:", err)
		return
	}
	hlog.Debugf("ws req: %s", req.String())
}

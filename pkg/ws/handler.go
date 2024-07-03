package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/websocket"
	"log"
	"open-copilot.dev/sidecar/pkg/completion"
	"sync"
)

var upgrader = websocket.HertzUpgrader{
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}
var wsMaxBufSize = 1024 * 1024 * 10
var wsPool = gopool.NewPool("ws", 1000, gopool.NewConfig())

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
	conn  *websocket.Conn
	buf   bytes.Buffer
	mutex sync.Mutex
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
	// 获取数据帧
	frame := ReadFrame(&h.buf)
	if frame == nil {
		return
	}

	// 解析 request
	var req = new(Request)
	err := json.Unmarshal(frame.Body, &req)
	if err != nil {
		hlog.Error("ws unmarshal:", err)
		return
	}
	hlog.Debugf("ws req: %s", req.String())

	// 处理 request
	if req.Method == "completion" {
		params := make([]completion.CompletionRequest, 0)
		err = json.Unmarshal(*req.Params, &params)
		if err != nil {
			hlog.Error("ws unmarshal:", err)
			h.sendError(req, err)
			return
		}
		if len(params) == 0 {
			hlog.Error("ws completion req is empty")
			h.sendError(req, errors.New("no completion request"))
			return
		}
		wsPool.Go(func() {
			result, err := completion.ProcessRequest(&params[0])
			if err != nil {
				hlog.Error("ws completion err:", err)
				h.sendError(req, err)
				return
			}
			h.sendResult(req, result)
			return
		})
	}
}

func (h *wsConnHandler) sendError(req *Request, err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	sendErr := h.conn.WriteJSON(&Response{
		Id:     req.Id,
		Result: nil,
		Error: &Error{
			Code:    -1,
			Message: err.Error(),
			Data:    nil,
		},
	})
	if sendErr != nil {
		hlog.Error("ws write:", sendErr)
	}
}

func (h *wsConnHandler) sendResult(req *Request, result any) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	sendErr := h.conn.WriteJSON(&Response{
		Id:     req.Id,
		Result: result,
	})
	if sendErr != nil {
		hlog.Error("ws write:", sendErr)
	}
}

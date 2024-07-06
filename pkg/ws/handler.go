package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	var wsRequest = new(Request)
	err := json.Unmarshal(frame.Body, &wsRequest)
	if err != nil {
		hlog.Error("ws unmarshal:", err)
		return
	}
	hlog.Debugf("wsRequest: %s", wsRequest.String())

	// 处理 request
	if wsRequest.Method == "completion" {
		completionParams := make([]completion.Request, 0)
		err = json.Unmarshal(*wsRequest.Params, &completionParams)
		if err != nil {
			hlog.Error("ws unmarshal:", err)
			h.sendError(wsRequest, err)
			return
		}
		if len(completionParams) == 0 {
			hlog.Error("ws completion params is empty")
			h.sendError(wsRequest, errors.New("no completion params"))
			return
		}
		wsPool.Go(func() {
			ctx := context.Background()
			completionResult, err := completion.ProcessRequest(ctx, &completionParams[0])
			if err != nil {
				hlog.Error("ws completion err:", err)
				h.sendError(wsRequest, err)
				return
			}
			h.send(&Response{
				Id:     wsRequest.Id,
				Result: completionResult,
			})
			return
		})
	}
}

func (h *wsConnHandler) sendError(wsRequest *Request, err error) {
	h.send(&Response{
		Id:     wsRequest.Id,
		Result: nil,
		Error: &Error{
			Code:    -1,
			Message: err.Error(),
			Data:    nil,
		},
	})
}

func (h *wsConnHandler) send(resp *Response) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	bodyBytes, err := json.Marshal(resp)
	if err != nil {
		hlog.Error("ws send marshal:", err)
		return
	}

	header := fmt.Sprintf("Content-Length: %d\r\n", len(bodyBytes))
	msg := header + "\r\n" + string(bodyBytes)

	sendErr := h.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if sendErr != nil {
		hlog.Error("ws send write:", sendErr)
	}
}

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
	"open-copilot.dev/sidecar/pkg/common"
	"open-copilot.dev/sidecar/pkg/completion"
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"sync"
)

var upgrader = websocket.HertzUpgrader{
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}
var wsMaxBufSize = 1024 * 1024 * 10
var wsPool = gopool.NewPool("ws", 1000, gopool.NewConfig())

func Handler(ctx context.Context, c *app.RequestContext) {
	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		wsHandler := &wsConnHandler{conn: conn}
		wsHandler.spin(ctx)
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "upgrade: %v", err)
		return
	}
}

// websocket connection handler
type wsConnHandler struct {
	conn      *websocket.Conn // 连接
	readBuf   bytes.Buffer    // 消息读取缓冲
	sendMutex sync.Mutex      // 消息发送锁
	reqCtxMap sync.Map        // 请求上下文
}

// ws message spin
func (h *wsConnHandler) spin(ctx context.Context) {
	for {
		_, message, err := h.conn.ReadMessage()
		if err != nil {
			hlog.CtxErrorf(ctx, "ws read: %v", err)
			break
		}
		hlog.CtxDebugf(ctx, "ws recv: %s", message)
		h.readBuf.Write(message)
		h.processRequests(ctx)
		if h.readBuf.Len() > wsMaxBufSize {
			// TODO 这边应该断开连接，因为直接清空缓冲池，可能导致消息错乱
			hlog.CtxErrorf(ctx, "ws readBuf is full, just clear it")
			h.readBuf.Reset()
		}
	}

	_ = h.conn.Close()
}

// 处理 ws 消息
func (h *wsConnHandler) processRequests(ctx context.Context) {
	for {
		// 获取数据帧
		frame := ReadFrame(&h.readBuf)
		if frame == nil {
			return
		}

		// 解析 request
		var wsRequest = new(Request)
		err := json.Unmarshal(frame.Body, &wsRequest)
		if err != nil {
			hlog.CtxErrorf(ctx, "ws unmarshal: %v", err)
			continue
		}
		hlog.CtxDebugf(ctx, "wsRequest: %s", wsRequest.String())

		// 请求上下文
		reqCtx := common.NewCancelableContext(context.WithValue(ctx, "X-Request-ID", wsRequest.Id))
		h.reqCtxMap.Store(wsRequest.Id, reqCtx)

		// 异步处理
		wsPool.Go(func() {
			defer h.reqCtxMap.Delete(wsRequest.Id)
			h.processRequest(reqCtx, wsRequest)
		})
	}
}

func (h *wsConnHandler) processRequest(ctx *common.CancelableContext, wsRequest *Request) {
	// 处理 request
	if wsRequest.Method == "completion" {
		h.processCompletionRequest(ctx, wsRequest)
	} else if wsRequest.Method == "$/cancelRequest" {
		h.processCancelRequest(ctx, wsRequest)
	}
}

func (h *wsConnHandler) processCompletionRequest(ctx *common.CancelableContext, wsRequest *Request) {
	completionParams := make([]domain.CompletionRequest, 0)
	err := json.Unmarshal(*wsRequest.Params, &completionParams)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws unmarshal: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	if len(completionParams) == 0 {
		hlog.CtxErrorf(ctx, "ws completion params is empty")
		h.sendError(wsRequest, errors.New("no completion params"))
		return
	}
	completionResult, err := completion.ProcessRequest(ctx, &completionParams[0])
	if err != nil {
		hlog.CtxErrorf(ctx, "ws completion err: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	h.send(&Response{
		Id:     wsRequest.Id,
		Result: completionResult,
	})
	return
}

func (h *wsConnHandler) processCancelRequest(ctx *common.CancelableContext, wsRequest *Request) {
	type cancelParam struct {
		ID string `json:"id"`
	}
	params := make([]cancelParam, 0)
	err := json.Unmarshal(*wsRequest.Params, &params)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws unmarshal: %v", err)
		return
	}
	if len(params) == 0 {
		hlog.CtxErrorf(ctx, "ws Cancel params is empty")
		return
	}

	val, ok := h.reqCtxMap.LoadAndDelete(params[0].ID)
	if !ok {
		return
	}
	willCancelReqCtx := val.(*common.CancelableContext)
	if willCancelReqCtx != nil {
		willCancelReqCtx.Cancel()
	}
}

func (h *wsConnHandler) sendError(wsRequest *Request, err error) {
	code := -1
	var errWithCode *common.ErrWithCode
	if errors.As(err, &errWithCode) {
		code = errWithCode.Code
	}

	h.send(&Response{
		Id:     wsRequest.Id,
		Result: nil,
		Error: &Error{
			Code:    code,
			Message: err.Error(),
			Data:    nil,
		},
	})
}

func (h *wsConnHandler) send(resp *Response) {
	if resp.Id == "" {
		hlog.Errorf("will not send because resp id is empty, resp: %+v", *resp)
		return
	}
	h.sendMutex.Lock()
	defer h.sendMutex.Unlock()

	bodyBytes, err := json.Marshal(resp)
	if err != nil {
		hlog.Error("ws send marshal:", err)
		return
	}

	header := fmt.Sprintf("Content-Length: %d\r\n", len(bodyBytes))
	msg := header + "\r\n" + string(bodyBytes) + "\r\n"

	sendErr := h.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if sendErr != nil {
		hlog.Error("ws send write:", sendErr)
	}
}

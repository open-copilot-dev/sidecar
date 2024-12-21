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
	"open-copilot.dev/sidecar/pkg/chat"
	chatDomain "open-copilot.dev/sidecar/pkg/chat/domain"
	"open-copilot.dev/sidecar/pkg/common"
	"open-copilot.dev/sidecar/pkg/completion"
	completionDomain "open-copilot.dev/sidecar/pkg/completion/domain"
	"sync"
)

var upgrader = websocket.HertzUpgrader{
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}

// 单个连接的消息缓冲池最大大小
var wsMaxBufSize = 1024 * 1024 * 10

// 消息处理协程池
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
			// 缓冲池满了，消息无法处理了，直接断开连接，客户端负责自动重试连接
			hlog.CtxErrorf(ctx, "ws readBuf is full, just close the connection")
			_ = h.conn.Close()
			break
		}
	}

	_ = h.conn.Close()
}

// -----------------------------------------------------------------

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
	} else if wsRequest.Method == "chat" {
		h.processChatRequest(ctx, wsRequest)
	} else if wsRequest.Method == "$/cancelRequest" {
		h.processCancelRequest(ctx, wsRequest)
	}
}

// 处理补全请求
func (h *wsConnHandler) processCompletionRequest(ctx *common.CancelableContext, wsRequest *Request) {
	wsParams := make([]completionDomain.CompletionRequest, 0)
	err := json.Unmarshal(*wsRequest.Params, &wsParams)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws unmarshal: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	if len(wsParams) == 0 {
		hlog.CtxErrorf(ctx, "ws completion params is empty")
		h.sendError(wsRequest, errors.New("no completion params"))
		return
	}
	completionRequest := &wsParams[0]
	completionResult, err := completion.ProcessRequest(ctx, completionRequest)
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

// 处理对话请求
func (h *wsConnHandler) processChatRequest(ctx *common.CancelableContext, wsRequest *Request) {
	wsParams := make([]chatDomain.ChatRequest, 0)
	err := json.Unmarshal(*wsRequest.Params, &wsParams)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws unmarshal: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	if len(wsParams) == 0 {
		hlog.CtxErrorf(ctx, "ws chat params is empty")
		h.sendError(wsRequest, errors.New("no chat params"))
		return
	}
	chatRequest := &wsParams[0]
	err = chat.ProcessRequest(ctx, chatRequest, func(streamResult *chatDomain.ChatStreamResult) {
		h.send(&Response{
			Id:     wsRequest.Id,
			Result: streamResult,
		})
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "ws chat err: %v", err)
		h.sendError(wsRequest, err)
		return
	}
}

// 处理取消请求
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

// -----------------------------------------------------------------
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

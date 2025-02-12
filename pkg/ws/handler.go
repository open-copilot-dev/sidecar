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
	"open-copilot.dev/sidecar/pkg/completion"
	"open-copilot.dev/sidecar/pkg/domain"
	"sync"
)

var upgrader = websocket.HertzUpgrader{
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}

// buf size of each connection
var wsMaxBufSize = 1024 * 1024 * 10

// connection process pool
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
	conn      *websocket.Conn // connection
	readBuf   bytes.Buffer    // message buffer
	sendMutex sync.Mutex      // message send lock
	reqCtxMap sync.Map        // request context map, key: requestID
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
		h.parseAndHandleRequests(ctx)
		if h.readBuf.Len() > wsMaxBufSize {
			// buffer is full, just close the connection, client will retry connect
			hlog.CtxErrorf(ctx, "ws readBuf is full, just close the connection")
			_ = h.conn.Close()
			break
		}
	}

	_ = h.conn.Close()
}

// -----------------------------------------------------------------

// parse ws request and handle it
func (h *wsConnHandler) parseAndHandleRequests(ctx context.Context) {
	for {
		// read ws data frame
		frame := ReadFrame(&h.readBuf)
		if frame == nil {
			return
		}

		// parse ws request
		var wsRequest = new(Request)
		err := json.Unmarshal(frame.Body, &wsRequest)
		if err != nil {
			hlog.CtxErrorf(ctx, "ws request unmarshal err: %v", err)
			continue
		}
		hlog.CtxDebugf(ctx, "recv ws request: %s", wsRequest.String())

		// create request context and store to context map
		reqCtx := domain.NewCancelableContext(context.WithValue(ctx, "X-Request-ID", wsRequest.Id))
		h.reqCtxMap.Store(wsRequest.Id, reqCtx)

		// async process request
		wsPool.Go(func() {
			defer h.reqCtxMap.Delete(wsRequest.Id)
			h.handleRequest(reqCtx, wsRequest)
		})
	}
}

// handleRequest handle single request
func (h *wsConnHandler) handleRequest(ctx *domain.CancelableContext, wsRequest *Request) {
	switch wsRequest.Method {
	case "completion":
		h.processCompletionRequest(ctx, wsRequest)
	case "chat":
		h.processChatRequest(ctx, wsRequest)
	case "chat/detail":
		h.processChatDetailRequest(ctx, wsRequest)
	case "chat/list":
		h.processChatListRequest(ctx, wsRequest)
	case "chat/deleteMessage":
		h.processChatMessageDeleteRequest(ctx, wsRequest)
	case "chat/delete":
		h.processChatDeleteRequest(ctx, wsRequest)
	case "chat/deleteAll":
		h.processChatDeleteAllRequest(ctx, wsRequest)
	case "$/cancelRequest":
		h.processCancelRequest(ctx, wsRequest)
	default:
		h.sendError(wsRequest, domain.ErrIllegal)
	}
}

// -----------------------------------------------------------------

// process completion request
func (h *wsConnHandler) processCompletionRequest(ctx *domain.CancelableContext, wsRequest *Request) {
	wsParams := make([]domain.CompletionRequest, 0)
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
		Jsonrpc: "2.0",
		Id:      wsRequest.Id,
		Result:  completionResult,
	})
	return
}

// process chat request
func (h *wsConnHandler) processChatRequest(ctx *domain.CancelableContext, wsRequest *Request) {
	wsParams := make([]domain.ChatRequest, 0)
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
	err = chat.ProcessRequest(ctx, chatRequest, func(streamResult *domain.ChatStreamResult) {
		h.send(&Response{
			Jsonrpc: "2.0",
			Id:      wsRequest.Id,
			Result:  streamResult,
		})
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "ws chat err: %v", err)
		h.sendError(wsRequest, err)
		return
	}
}
func (h *wsConnHandler) processChatDetailRequest(ctx *domain.CancelableContext, wsRequest *Request) {
	wsParams := make([]string, 0)
	err := json.Unmarshal(*wsRequest.Params, &wsParams)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws unmarshal: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	if len(wsParams) == 0 {
		hlog.CtxErrorf(ctx, "ws chat detail params is empty")
		h.sendError(wsRequest, errors.New("no chat detail params"))
		return
	}
	chatID := wsParams[0]
	chatDetail, err := chat.ProcessDetailRequest(ctx, chatID)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws chat detail err: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	h.send(&Response{
		Jsonrpc: "2.0",
		Id:      wsRequest.Id,
		Result:  chatDetail,
	})
}
func (h *wsConnHandler) processChatListRequest(ctx *domain.CancelableContext, wsRequest *Request) {
	chats, err := chat.ProcessListRequest(ctx)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws chat list err: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	h.send(&Response{
		Jsonrpc: "2.0",
		Id:      wsRequest.Id,
		Result:  chats,
	})
}
func (h *wsConnHandler) processChatMessageDeleteRequest(ctx *domain.CancelableContext, wsRequest *Request) {
	wsParams := make([]string, 0)
	err := json.Unmarshal(*wsRequest.Params, &wsParams)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws unmarshal: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	if len(wsParams) < 2 {
		hlog.CtxErrorf(ctx, "ws chat delete message params missing")
		h.sendError(wsRequest, errors.New("missing chat delete message params"))
		return
	}
	chatID := wsParams[0]
	messageID := wsParams[1]
	err = chat.ProcessDeleteMessageRequest(ctx, chatID, messageID)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws chat message delete err: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	h.send(&Response{
		Jsonrpc: "2.0",
		Id:      wsRequest.Id,
		Result:  true,
	})
}
func (h *wsConnHandler) processChatDeleteRequest(ctx *domain.CancelableContext, wsRequest *Request) {
	wsParams := make([]string, 0)
	err := json.Unmarshal(*wsRequest.Params, &wsParams)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws unmarshal: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	if len(wsParams) == 0 {
		hlog.CtxErrorf(ctx, "ws chat delete params is empty")
		h.sendError(wsRequest, errors.New("no chat delete params"))
		return
	}
	chatID := wsParams[0]
	err = chat.ProcessDeleteRequest(ctx, chatID)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws chat delete err: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	h.send(&Response{
		Jsonrpc: "2.0",
		Id:      wsRequest.Id,
		Result:  true,
	})
}
func (h *wsConnHandler) processChatDeleteAllRequest(ctx *domain.CancelableContext, wsRequest *Request) {
	err := chat.ProcessDeleteAllRequest(ctx)
	if err != nil {
		hlog.CtxErrorf(ctx, "ws chat delete all err: %v", err)
		h.sendError(wsRequest, err)
		return
	}
	h.send(&Response{
		Jsonrpc: "2.0",
		Id:      wsRequest.Id,
		Result:  true,
	})
}

// -----------------------------------------------------------------

// process cancel request
func (h *wsConnHandler) processCancelRequest(ctx *domain.CancelableContext, wsRequest *Request) {
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
	willCancelReqCtx := val.(*domain.CancelableContext)
	if willCancelReqCtx != nil {
		willCancelReqCtx.Cancel()
	}
}

// -----------------------------------------------------------------

// send err resp
func (h *wsConnHandler) sendError(wsRequest *Request, err error) {
	code := -1
	var errWithCode *domain.ErrWithCode
	if errors.As(err, &errWithCode) {
		code = errWithCode.Code
	}

	h.send(&Response{
		Jsonrpc: "2.0",
		Id:      wsRequest.Id,
		Result:  nil,
		Error: &Error{
			Code:    code,
			Message: err.Error(),
			Data:    nil,
		},
	})
}

// send resp
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

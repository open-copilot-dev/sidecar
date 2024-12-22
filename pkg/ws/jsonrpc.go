package ws

import "encoding/json"

//----------------------------------------
// jsonrpc request

type Request struct {
	Jsonrpc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params"`
	Id      string           `json:"id"`
}

func (req *Request) String() string {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

//----------------------------------------
// jsonrpc response

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result  any    `json:"result"`
	Error   any    `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

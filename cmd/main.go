package main

import (
	"flag"
	"github.com/cloudwego/hertz/pkg/app/server"
	"open-copilot.dev/sidecar/pkg/ws"
)

var addr = flag.String("addr", "localhost:30999", "http service address")
var debug = flag.Bool("debug", false, "debug mode")

func main() {
	initLog(debug)

	flag.Parse()
	h := server.Default(server.WithHostPorts(*addr))
	h.NoHijackConnPool = true
	h.GET("/ws", ws.Handler)
	h.Spin()
}

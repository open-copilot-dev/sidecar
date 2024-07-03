package ws

import (
	"bytes"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strconv"
	"strings"
)

/*
兼容 lsp 规定的消息格式 https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/

	Content-Length: ...\r\n
	\r\n
	{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "textDocument/completion",
		"params": {
			...
		}
	}
*/

type Frame struct {
	Header *FrameHeader
	Body   []byte
}

type FrameHeader struct {
	ContentLength int
}

func ReadFrame(buf *bytes.Buffer) *Frame {
	bs := buf.Bytes()

	header, readLen := ReadFrameHeader(bs)
	if header == nil {
		// 未能成功读取 header
		hlog.Warn("failed to read frame header")
		return nil
	}
	if header.ContentLength <= 0 {
		// 不合法的 header 跳过
		hlog.Warnf("invalid frame header, skip %d bytes", readLen)
		buf.Next(readLen)
		return nil
	}
	if len(bs) < readLen+header.ContentLength {
		// Body 部分还不够，先不读取
		return nil
	}
	buf.Next(readLen)
	body := buf.Next(header.ContentLength)
	return &Frame{
		Header: header,
		Body:   body,
	}
}

func ReadFrameHeader(bs []byte) (*FrameHeader, int) {
	var totalLen = len(bs)
	var contentLength *int = nil

	for {
		idx := bytes.IndexByte(bs, '\n')
		if idx == -1 {
			break
		}
		line := string(bs[:idx+1])
		bs = bs[idx+1:]
		if line == "" || line == "\n" || line == "\r\n" {
			// 空行，代表header结束了
			break
		}
		colon := strings.IndexRune(line, ':')
		if colon < 0 {
			// 此行在header中是非法的，跳过
			hlog.Warn("invalid header line: " + line)
			continue
		}
		name, value := line[:colon], strings.TrimSpace(line[colon+1:])
		switch name {
		case "Content-Length":
			length, err := strconv.Atoi(value)
			if err != nil {
				hlog.Warn("invalid Content-Length: " + value)
				continue
			} else {
				contentLength = &length
			}
		default:
			// ignoring unknown headers
		}
	}

	if contentLength == nil {
		return nil, -1
	}
	readLen := totalLen - len(bs)
	return &FrameHeader{ContentLength: *contentLength}, readLen
}

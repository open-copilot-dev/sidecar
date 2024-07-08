package ws

import (
	"bytes"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"testing"
)

func TestReadFrameHeader(t *testing.T) {
	strHeader := "\r\nContent-Length: 3\r\n"
	strBody := "abc"
	bs := []byte(strHeader + "\r\n" + strBody)
	header, readLen := ReadFrameHeader(bs)
	assert.Assert(t, header.ContentLength == 3)
	assert.Assert(t, readLen == len(strHeader)+2)

	strHeader = "Content-Length: 3\n"
	strBody = "abc"
	bs = []byte(strHeader + "\n" + strBody)
	header, readLen = ReadFrameHeader(bs)
	assert.Assert(t, header.ContentLength == 3)
	assert.Assert(t, readLen == len(strHeader)+1)

	strHeader = "Content-Length: 3\r\nContent-Type: application/vscode-jsonrpc; charset=utf-8\r\n"
	strBody = "abc"
	bs = []byte(strHeader + "\r\n" + strBody)
	header, readLen = ReadFrameHeader(bs)
	assert.Assert(t, header.ContentLength == 3)
	assert.Assert(t, readLen == len(strHeader)+2)

	strHeader = "Content-Type: application/vscode-jsonrpc; charset=utf-8\r\nContent-Length: 3\r\n"
	strBody = "abc"
	bs = []byte(strHeader + "\r\n" + strBody)
	header, readLen = ReadFrameHeader(bs)
	assert.Assert(t, header.ContentLength == 3)
	assert.Assert(t, readLen == len(strHeader)+2)

	strHeader = "xxx\r\n"
	strBody = "abc"
	bs = []byte(strHeader + "\r\n" + strBody)
	header, readLen = ReadFrameHeader(bs)
	assert.Assert(t, header == nil)
}

func TestReadFrame(t *testing.T) {
	strHeader := "Content-Length: 3\r\n"
	strBody := "abc"
	bs := []byte(strHeader + "\r\n" + strBody)
	frame := ReadFrame(bytes.NewBuffer(bs))
	assert.Assert(t, frame.Header.ContentLength == 3)
	assert.Assert(t, string(frame.Body) == strBody)

	strHeader = "Content-Length: 3\n"
	strBody = "abc"
	bs = []byte(strHeader + "\n" + strBody)
	frame = ReadFrame(bytes.NewBuffer(bs))
	assert.Assert(t, frame.Header.ContentLength == 3)
	assert.Assert(t, string(frame.Body) == strBody)

	strHeader = "Content-Length: 3\r\nContent-Type: application/vscode-jsonrpc; charset=utf-8\r\n"
	strBody = "abc"
	bs = []byte(strHeader + "\r\n" + strBody)
	frame = ReadFrame(bytes.NewBuffer(bs))
	assert.Assert(t, frame.Header.ContentLength == 3)
	assert.Assert(t, string(frame.Body) == strBody)

	strHeader = "Content-Type: application/vscode-jsonrpc; charset=utf-8\r\nContent-Length: 3\r\n"
	strBody = "abc"
	bs = []byte(strHeader + "\r\n" + strBody)
	frame = ReadFrame(bytes.NewBuffer(bs))
	assert.Assert(t, frame.Header.ContentLength == 3)
	assert.Assert(t, string(frame.Body) == strBody)

	strHeader = "Content-Type: application/vscode-jsonrpc; charset=utf-8\r\nContent-Length: 3\r\n"
	strBody = "ab"
	bs = []byte(strHeader + "\r\n" + strBody)
	frame = ReadFrame(bytes.NewBuffer(bs))
	assert.Assert(t, frame == nil)

	strHeader = "Content-Type: application/vscode-jsonrpc; charset=utf-8\r\nContent-Length: 3\r\n"
	strBody = "abcdef"
	bs = []byte(strHeader + "\r\n" + strBody)
	buffer := bytes.NewBuffer(bs)
	frame = ReadFrame(buffer)
	assert.Assert(t, string(frame.Body) == "abc")
	assert.Assert(t, string(buffer.Bytes()) == "def")
}

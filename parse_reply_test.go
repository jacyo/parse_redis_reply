package main

import (
	"testing"
	"strings"
)

func TestReplyReader(t *testing.T) {
	rr := NewReplyReader(strings.NewReader("+OK\r\n"))
	t.Log(rr.ParseReply())

	rr = NewReplyReader(strings.NewReader("-Error\r\n"))
	t.Log(rr.ParseReply())

	rr = NewReplyReader(strings.NewReader(":1000\r\n"))
	t.Log(rr.ParseReply())

	rr = NewReplyReader(strings.NewReader("$6\r\nfoobar\r\n"))
	t.Log(rr.ParseReply())

	rr = NewReplyReader(strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
	t.Log(rr.ParseReply())

	rr = NewReplyReader(strings.NewReader("*3\r\n$5\r\nhello\r\n$-1\r\n$5\r\nworld\r\n"))
	t.Log(rr.ParseReply())
}
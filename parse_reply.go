package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const (
	SimpleReply = '+'
	ErrorReply  = '-'
	IntReply    = ':'
	BulkReply   = '$'
	ArraysReply = '*'
)

var RedisNil = errors.New("redis: nil")

type ReplyReader struct {
	rd *bufio.Reader
}

func NewReplyReader(rd io.Reader) *ReplyReader {
	return &ReplyReader{
		rd: bufio.NewReader(rd),
	}
}

func (r *ReplyReader) ReadLine() ([]byte, error) {
	line, isPrefix, err := r.rd.ReadLine()
	if err != nil {
		return nil, err
	}
	if isPrefix {
		return nil, bufio.ErrBufferFull
	}
	if len(line) == 0 {
		return nil, fmt.Errorf("redis: line is empty")
	}
	if isNil(line) {
		return nil, RedisNil
	}
	return line, nil
}

func (r *ReplyReader) ParseReply() (interface{}, error) {
	line, err := r.ReadLine()
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case SimpleReply, ErrorReply:
		return string(line[1:]), nil
	case IntReply:
		return strconv.ParseInt(string(line[1:]), 10, 64)
	case BulkReply:
		return r.ParseBulkReply(line)
	case ArraysReply:
		return r.ParseArraysReply(line)
	}
	return nil, fmt.Errorf("redis: unsupported protocols!")
}

func (r *ReplyReader) ParseBulkReply(line []byte) (string, error) {
	if isNil(line) {
		return "", RedisNil
	}

	length, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return "", err
	}

	b := make([]byte, length+2)
	_, err = io.ReadFull(r.rd, b)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (r *ReplyReader) ParseArraysReply(line []byte) ([]string, error) {
	if isNil(line) {
		return nil, RedisNil
	}

	length, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, err
	}

	replySlice := make([]string, 0, length)
	for i := 0; i < length; i++ {
		line2, err := r.ReadLine()
		if err == RedisNil {
			replySlice = append(replySlice, "")
			continue
		}
		if err != nil {
			return nil, err
		}

		if line2[0] != BulkReply {
			return nil, fmt.Errorf("redis: parse error!")
		}

		s, err := r.ParseBulkReply(line2)
		if err == RedisNil {
			replySlice = append(replySlice, "")
			continue
		}
		if err != nil {
			return nil, err
		}

		replySlice = append(replySlice, s)
	}
	return replySlice, nil
}

func isNil(b []byte) bool {
	return len(b) == 3 &&
		(b[0] == BulkReply || b[0] == ArraysReply) &&
		b[1] == '-' && b[2] == '1'
}

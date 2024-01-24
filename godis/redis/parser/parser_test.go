package parser

import (
	"bytes"
	"io"
	"testing"

	"example.com/redis/interface/redis"
	"example.com/redis/lib/utils"
	"example.com/redis/redis/protocol"
)

func TestParseStream(t *testing.T) {
	replies := []redis.Reply{
		protocol.NewIntReply(1),
		protocol.NewSimpleReply("OK"),
		protocol.NewStandardErrorReply("Err unknown"),
		protocol.NewBulkReply([]byte("a\r\nb")),
		protocol.NewBulkArrayReply([][]byte{
			[]byte("a"),
			[]byte("\r\n"),
		}),
	}

	reqs := bytes.Buffer{}
	for _, re := range replies {
		reqs.Write(re.ToBytes())
	}
	reqs.Write([]byte("set a a" + protocol.CRLF))

	expected := make([]redis.Reply, len(replies))
	copy(expected, replies)
	expected = append(expected, protocol.NewBulkArrayReply([][]byte{
		[]byte("set"), []byte("a"), []byte("a"),
	}))

	ch := ParseStream(bytes.NewReader(reqs.Bytes()))
	i := 0
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF {
				return
			} else {
				t.Error(payload.Err)
				return
			}
		} else if payload.Data == nil {
			t.Error("Empty Data")
			return
		}
		exp := expected[i]
		i++
		if !utils.BytesEqual(exp.ToBytes(), payload.Data.ToBytes()) {
			t.Error("parse failed.\nexpected: " + string(exp.ToBytes()) + "get: " + string(payload.Data.ToBytes()))
			return
		}
	}
}

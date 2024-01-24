package protocol

import (
	"bytes"
	"strconv"

	"example.com/redis/interface/redis"
)

var CRLF = "\r\n"

/* ------ Simple String Reply -------*/

type SimpleReply struct {
	Message string
}

func NewSimpleReply(message string) *SimpleReply {
	return &SimpleReply{
		Message: message,
	}
}

func (r *SimpleReply) ToBytes() []byte {
	return []byte("+" + r.Message + CRLF)
}

func IsOKReply(r redis.Reply) bool {
	return string(r.ToBytes()) == "+OK\r\n"
}

/* ------ Integer Reply -------- */

type IntReply struct {
	Code int64
}

func NewIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

func (r *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(r.Code, 10) + CRLF)
}

/* ------- Error Reply -------- */

type StandardErrorReply struct {
	Message string
}

func NewStandardErrorReply(message string) *StandardErrorReply {
	return &StandardErrorReply{
		Message: message,
	}
}

func IsErrorReply(r redis.Reply) bool {
	return r.ToBytes()[0] == '-'
}

func (r *StandardErrorReply) Error() string {
	return r.Message
}

func (r *StandardErrorReply) ToBytes() []byte {
	return []byte("-" + r.Message + CRLF)
}

/* ------ Bulk String Reply ------ */

// bulk string stores a binary-safe string.
type BulkReply struct {
	Arg []byte
}

func NewBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

func (r *BulkReply) ToBytes() []byte {
	if r.Arg == nil {
		return []byte("$-1" + CRLF)
	}
	return []byte("$" + strconv.Itoa(len(r.Arg)) + CRLF + string(r.Arg) + CRLF)
}

/* ------- Bulk Array String ---- */

type BulkArrayReply struct {
	Args [][]byte
}

func NewBulkArrayReply(args [][]byte) *BulkArrayReply {
	return &BulkArrayReply{
		Args: args,
	}
}

func (r *BulkArrayReply) ToBytes() []byte {
	argsLen := len(r.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argsLen) + CRLF)
	for _, arg := range r.Args {
		if arg == nil {
			buf.WriteString("$-1" + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

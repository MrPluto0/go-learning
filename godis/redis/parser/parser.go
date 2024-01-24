package parser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"

	"example.com/redis/interface/redis"
	"example.com/redis/redis/protocol"
)

type Payload struct {
	Data redis.Reply
	Err  error
}

func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse(reader, ch)
	return ch
}

func parse(rawReader io.Reader, ch chan<- *Payload) {
	reader := bufio.NewReader(rawReader)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			ch <- &Payload{Err: err}
			close(ch)
			return
		}
		length := len(line)
		if length <= 2 || line[length-2] != '\r' {
			continue
		}
		line = bytes.TrimSuffix(line, []byte(protocol.CRLF))
		switch line[0] {
		case '+':
			ch <- &Payload{
				Data: protocol.NewSimpleReply(string(line[1:])),
			}
		case '-':
			ch <- &Payload{
				Data: protocol.NewStandardErrorReply(string(line[1:])),
			}
		case ':':
			code, err := strconv.ParseInt(string(line[1:]), 10, 64)
			if err != nil {
				protocolError(ch, "illegal number: "+string(line[1:]))
				continue
			}
			ch <- &Payload{
				Data: protocol.NewIntReply(code),
			}
		case '$':
			err := parseBulkString(ch, line, reader)
			if err != nil {
				ch <- &Payload{Err: err}
				close(ch)
				return
			}
		case '*':
			err := parseBulkArray(ch, line, reader)
			if err != nil {
				ch <- &Payload{Err: err}
				close(ch)
				return
			}
		default:
			args := bytes.Split(line, []byte(" "))
			ch <- &Payload{
				Data: protocol.NewBulkArrayReply(args),
			}
		}
	}
}

func parseBulkString(ch chan<- *Payload, header []byte, reader *bufio.Reader) error {
	length, err := strconv.ParseInt(string(header[1:]), 10, 64)
	if err != nil || length < -1 {
		protocolError(ch, "illegal bulk string length: "+string(header))
		return nil
	} else if length == -1 {
		ch <- &Payload{
			Data: protocol.NewBulkReply(nil),
		}
		return nil
	} else {
		body := make([]byte, length+2)
		_, err := io.ReadFull(reader, body)
		if err != nil {
			return err
		}
		ch <- &Payload{
			Data: protocol.NewBulkReply(body[:len(body)-2]),
		}
		return nil
	}
}

func parseBulkArray(ch chan<- *Payload, header []byte, reader *bufio.Reader) error {
	argsLen, err := strconv.ParseInt(string(header[1:]), 10, 64)
	if err != nil {
		protocolError(ch, "illegal bulk array length: "+string(header))
		return nil
	}
	lines := make([][]byte, 0, argsLen)
	for i := int64(0); i < argsLen; i++ {
		subHeader, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		length := len(subHeader)
		if length < 4 || subHeader[length-2] != '\r' || subHeader[0] != '$' {
			protocolError(ch, "illegal bulk array element header: "+string(subHeader))
			break
		}
		argLen, err := strconv.ParseInt(string(subHeader[1:length-2]), 10, 64)
		if err != nil || argLen < -1 {
			protocolError(ch, "illegal bulk array element length: "+string(subHeader))
			break
		} else if argLen == -1 {
			lines = append(lines, nil)
		} else {
			body := make([]byte, argLen+2)
			_, err := io.ReadFull(reader, body)
			if err != nil {
				return err
			}
			lines = append(lines, body[:len(body)-2])
		}
	}
	ch <- &Payload{
		Data: protocol.NewBulkArrayReply(lines),
	}
	return nil
}

func protocolError(ch chan<- *Payload, msg string) {
	err := errors.New("protocol error: " + msg)
	ch <- &Payload{Err: err}
}

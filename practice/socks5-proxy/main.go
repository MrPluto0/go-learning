package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

const socks5Ver = 0x05
const cmdBind = 0x01
const atypIPV4 = 0x01
const atypeHOST = 0x03
const atypeIPV6 = 0x04

func main() {
	server, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		panic(err)
	}
	log.Printf("Listen to 127.0.0.1:8081...")

	for {
		client, err := server.Accept()
		if err != nil {
			log.Printf("Accept failed %v", err)
			continue
		}
		go process(client)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	err := auth(reader, conn)
	if err != nil {
		log.Printf("client %v auth failed:%v\n", conn.RemoteAddr(), err)
	}

	err = connect(reader, conn)
	if err != nil {
		log.Printf("client %v connect failed:%v\n", conn.RemoteAddr(), err)
	}
}

// type ProtocolVersion struct {
//     VER uint8
//     NMETHODS uint8
//     METHODS []uint8
// }

func auth(reader *bufio.Reader, conn net.Conn) error {
	/*
		socks5收到客户端认证消息格式：
		+----+----------+----------+
		|VER | NMETHODS | METHODS  |
		+----+----------+----------+
		| 1  |    1     | 1 to 255 |
		+----+----------+----------+
		VER: 协议版本，socks5为0x05
		NMETHODS:METHODS部分的长度；
		METHODS: METHODS是客户端支持的认证方式列表，每个方法占1字节。当前的定义是：
		  - 0x01 GSSAPI
		  - 0x02 用户名、密码认证
		  - 0x03 - 0x7F由IANA分配（保留）
		  - 0x80 - 0xFE为私人方法保留
		  - 0xFF 无可接受的方法

		sockes5返回消息格式：
		+----+--------+
		|VER | METHOD |
		+----+--------+
		| 1  |   1    |
		+----+--------+
	*/

	ver, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read version failed: %w", err)
	}
	if ver != socks5Ver {
		return fmt.Errorf("not supported version: %v", ver)
	}

	methodSize, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read method size failed: %w", err)
	}
	method := make([]byte, methodSize)
	_, err = io.ReadFull(reader, method)
	if err != nil {
		return fmt.Errorf("read method body failed: %w", err)
	}

	log.Println("version", ver, "method", method)

	// method=0x00 代表不需要认证
	_, err = conn.Write([]byte{socks5Ver, 0x00})
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	return nil
}

func connect(reader *bufio.Reader, conn net.Conn) error {
	/*
			socks5收到客户端请求格式：
			+----+-----+-------+------+----------+----------+
			|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
			+----+-----+-------+------+----------+----------+
			| 1  |  1  | X'00' |  1   | Variable |    2     |
			+----+-----+-------+------+----------+----------+
			VER 版本号，socks5的值为0x05
			CMD 0x01表示CONNECT请求
			RSV 保留字段，值为0x00
			ATYP 目标地址类型，DST.ADDR的数据对应这个字段的类型。
				- 0x01 IPv4地址，DST.ADDR部分4字节长度
				- 0x03 域名，DST.ADDR部分第一个字节为域名长度，DST.ADDR剩余的内容为域名，没有\0结尾。
				- 0x04 IPv6地址，16个字节长度。
			DST.ADDR 一个可变长度的值
			DST.PORT 目标端口，固定2个字节

			socks5发送客户端消息格式：
		  +----+-----+-------+------+----------+----------+
		  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
		  +----+-----+-------+------+----------+----------+
		  | 1  |  1  | X'00' |  1   | Variable |    2     |
		  +----+-----+-------+------+----------+----------+
		  VER socks版本，这里为0x05
		  REP Relay field,内容取值如下 X’00’ succeeded
		  RSV 保留字段
		  ATYPE 地址类型
		  BND.ADDR 服务绑定的地址
		  BND.PORT 服务绑定的端口DST.PORT
	*/

	buf := make([]byte, 4)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return fmt.Errorf("read header failed:%w", err)
	}

	ver, cmd, atyp := buf[0], buf[1], buf[3]
	if ver != socks5Ver {
		return fmt.Errorf("not supported version:%v", ver)
	}
	if cmd != cmdBind {
		return fmt.Errorf("not supported cmd:%v", cmd)
	}

	var addr string
	switch atyp {
	case atypIPV4:
		_, err = io.ReadFull(reader, buf) // 读取四个字节
		if err != nil {
			return fmt.Errorf("read atyp failed:%w", err)
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case atypeHOST:
		hostSize, err := reader.ReadByte()
		if err != nil {
			return fmt.Errorf("read hostsize failed:%w", err)
		}
		host := make([]byte, hostSize)
		_, err = io.ReadFull(reader, host)
		if err != nil {
			return fmt.Errorf("read host failed: %w", err)
		}
		addr = string(host)
	case atypeIPV6:
		return errors.New("IPV6: not supported yet")
	default:
		return errors.New("invalid atyp")
	}

	_, err = io.ReadFull(reader, buf[:2])
	if err != nil {
		return fmt.Errorf("read port failed:%w", err)
	}
	port := binary.BigEndian.Uint16(buf[:2])

	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	dest, err := net.Dial("tcp", fmt.Sprintf("%v:%v", addr, port))
	if err != nil {
		return fmt.Errorf("dial dst failed:%w", err)
	}
	defer dest.Close()
	log.Println("dial", addr, port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_, _ = io.Copy(dest, conn)
		cancel()
	}()

	go func() {
		_, _ = io.Copy(conn, dest)
		cancel()
	}()

	<-ctx.Done()
	return nil
}

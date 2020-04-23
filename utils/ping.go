package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"
)

//ICMP 定义icmp结构体
type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

//Ret 返回结果
type Ret struct {
	Success       bool
	RetBytesLenth int
	TTL           byte
	Et            int
}

var (
	icmp  ICMP
	laddr = net.IPAddr{IP: net.ParseIP("ip")}
)

//CheckSum 检查校验和
func CheckSum(data []byte) uint16 {
	var sum uint32
	var length = len(data)
	var index int

	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length == 1 {
		sum += uint32(data[index])
	}
	sum = uint32(sum>>16) + uint32(sum)
	sum = uint32(sum>>16) + uint32(sum)
	return uint16(^sum)
}

//Ping ping目标地址
func Ping(desIP string, size int, timeout int64) Ret {
	var ret Ret

	// Dial icmp
	conn, err := net.DialTimeout("ip:icmp", desIP, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	// 退出时关闭连接
	defer conn.Close()

	//icmp头部填充
	icmp.Type = 8
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 1
	icmp.SequenceNum = 1

	// 定义发送buffer
	var buffer bytes.Buffer
	// 写入二进制buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	// 定义data
	data := make([]byte, size)

	buffer.Write(data)
	data = buffer.Bytes()

	icmp.SequenceNum = uint16(1)
	// 检验和设为0
	data[2] = byte(0)
	data[3] = byte(0)

	data[6] = byte(icmp.SequenceNum >> 8)
	data[7] = byte(icmp.SequenceNum)

	icmp.Checksum = CheckSum(data)
	data[2] = byte(icmp.Checksum >> 8)
	data[3] = byte(icmp.Checksum)

	// 开始时间
	t1 := time.Now()
	conn.SetDeadline(t1.Add(time.Duration(time.Duration(timeout) * time.Millisecond)))
	n, err := conn.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	// 定义读取buf
	buf := make([]byte, 65535)
	// 读取返回数据到buf
	n, err = conn.Read(buf)

	// 请求超时增加失败次数
	if err != nil {
		ret.Success = false
		return ret
	}

	et := int(time.Since(t1) / 1000000)

	ret.TTL = buf[8]
	ret.RetBytesLenth = len(buf[28:n])
	ret.Et = et

	ret.Success = true
	return ret
}

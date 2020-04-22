package main

import (
	"flag"
	"fmt"
	"net"
)

//ICMP 定义icmp结构体
type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

var (
	icmp    ICMP
	laddr   = net.IPAddr{IP: net.ParseIP("ip")}
	num     int
	timeout int64
	size    int
	stop    bool
)

func main() {
	fmt.Println("helloworld")
}

//ParseArgs 解析参数
func ParseArgs() {
	flag.Int64Var(&timeout, "w", 1000, "等待每次回复的超时时间（毫秒）")
	flag.IntVar(&num, "n", 4, "要发送的请求数")
	flag.IntVar(&size, "1", 32, "要发送缓冲区大小")
	flag.BoolVar(&stop, "t", false, "Ping 指定的主机，直到停止")

	flag.Parse()
}

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
	sum = uint16(sum>>16) + uint16(sum)
	sum = uint16(sum>>16) + uint16(sum)
	return uint16(^sum)
}

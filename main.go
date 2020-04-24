package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/sersoong/go-ping/utils"
)

var (
	num     int
	timeout int64
	size    int
	stop    bool
)

func main() {
	// 解析传入参数
	ParseArgs()
	args := os.Args
	desIP := args[len(args)-1]
	// 如果参数数量小于2则提示help信息
	if len(args) < 2 {
		Usage()
	}

	fmt.Printf("\n正在 ping %s 具有 %d 字节的数据:\n", desIP, size)

	defer func() {
		fmt.Println("recover start...")
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("recover end")
	}()
	var SuccessTimes int
	var FailTimes int
	var totalTime int
	var minTime = int(math.MaxInt32)
	var maxTime int
	for i := 0; i < num; i++ {
		ret := utils.Ping(desIP, size, timeout)
		if ret.Success {
			SuccessTimes++
			if minTime > ret.Et {
				minTime = ret.Et
			}
			if maxTime < ret.Et {
				maxTime = ret.Et
			}
			fmt.Printf("来自 %s 的回复: 字节=%d 时间=%dms TTL=%d\n", desIP, ret.RetBytesLenth, ret.Et, ret.TTL)
			totalTime += ret.Et
		} else {
			fmt.Println("请求超时。")
			FailTimes++
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\n%s 的 Ping 统计信息:\n", desIP)
	fmt.Printf("    数据包: 已发送 = %d，已接收 = %d，丢失 = %d (%.2f%% 丢失)，\n", SuccessTimes+FailTimes, SuccessTimes, FailTimes, float64(FailTimes*100)/float64(SuccessTimes+FailTimes))
	if maxTime != 0 && minTime != int(math.MaxInt32) {
		fmt.Printf("往返行程的估计时间(以毫秒为单位):\n")
		fmt.Printf("    最短 = %dms，最长 = %dms，平均 = %dms\n", minTime, maxTime, totalTime/SuccessTimes)
	}
}

//ParseArgs 解析参数
func ParseArgs() {
	flag.Int64Var(&timeout, "w", 1000, "等待每次回复的超时时间（毫秒）")
	flag.IntVar(&num, "n", 4, "要发送的请求数")
	flag.IntVar(&size, "l", 32, "要发送缓冲区大小")
	flag.BoolVar(&stop, "t", false, "Ping 指定的主机，直到停止")

	flag.Parse()
}

//Usage 打印使用方法
func Usage() {
	argNum := len(os.Args)
	if argNum < 2 {
		fmt.Println(`
用法: ping [-t] [-a] [-n count] [-l size] [-f] [-i TTL] [-v TOS]
            [-r count] [-s count] [[-j host-list] | [-k host-list]]
            [-w timeout] [-R] [-S srcaddr] [-c compartment] [-p]
            [-4] [-6] target_name
选项:
    -t             Ping 指定的主机，直到停止。
                   若要查看统计信息并继续操作，请键入 Ctrl+Break；
                   若要停止，请键入 Ctrl+C。
    -a             将地址解析为主机名。
    -n count       要发送的回显请求数。
    -l size        发送缓冲区大小。
    -f             在数据包中设置“不分段”标记(仅适用于 IPv4)。
    -i TTL         生存时间。
    -v TOS         服务类型(仅适用于 IPv4。该设置已被弃用，
                   对 IP 标头中的服务类型字段没有任何
                   影响)。
    -r count       记录计数跃点的路由(仅适用于 IPv4)。
    -s count       计数跃点的时间戳(仅适用于 IPv4)。
    -j host-list   与主机列表一起使用的松散源路由(仅适用于 IPv4)。
    -k host-list    与主机列表一起使用的严格源路由(仅适用于 IPv4)。
    -w timeout     等待每次回复的超时时间(毫秒)。
    -R             同样使用路由标头测试反向路由(仅适用于 IPv6)。
                   根据 RFC 5095，已弃用此路由标头。
                   如果使用此标头，某些系统可能丢弃
                   回显请求。
    -S srcaddr     要使用的源地址。
    -c compartment 路由隔离舱标识符。
    -p             Ping Hyper-V 网络虚拟化提供程序地址。
    -4             强制使用 IPv4。
    -6             强制使用 IPv6。
`)
	}
}

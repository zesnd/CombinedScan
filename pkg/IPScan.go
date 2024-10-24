package pkg

import (
	"fmt"
	"golang.org/x/exp/rand"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ListenAddr 默认监听IPV4端口
var ListenAddr = "0.0.0.0"

// Ping 执行ICMP Ping操作
func Ping(addr string) (*net.IPAddr, time.Duration, error) {
	// 监听ICMP回包
	conn, err := icmp.ListenPacket("ip4:icmp", ListenAddr)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close()

	// 进行DNS解析，并返回真实IP
	dst, err := net.ResolveIPAddr("ip4", addr)
	if err != nil {
		return nil, 0, err
	}

	// 生成随机数据
	data := make([]byte, 32)
	_, err = rand.Read(data)
	if err != nil {
		return dst, 0, err
	}

	// 构建ICMP包
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: data,
		},
	}

	b, err := m.Marshal(nil)
	if err != nil {
		return dst, 0, err
	}

	// 发送ICMP包
	start := time.Now()
	n, err := conn.WriteTo(b, dst)
	if err != nil {
		return dst, time.Since(start), err
	} else if n != len(b) {
		return dst, 0, fmt.Errorf("got %v; want %v", n, len(b))
	}

	// 接收回包
	reply := make([]byte, 1024)
	err = conn.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		return dst, time.Since(start), err
	}
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		return dst, time.Since(start), err
	}

	// 判断处理
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), reply[:n])
	if err != nil {
		return dst, time.Since(start), err
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return dst, time.Since(start), nil
	default:
		return dst, 0, fmt.Errorf("got %v; want %v", rm, peer)
	}
}

// IPScan 执行IP扫描
func IPScan(targetIP string) ([]string, error) {
	// 解析命令行参数，获取目标 IP 地址
	if targetIP == "" {
		return nil, fmt.Errorf("请输入要扫描的目标IP地址")
	}

	// 创建通道用于接收结果
	resultChan := make(chan string)

	// 用于等待所有 goroutine 完成
	var wg sync.WaitGroup

	// 解析IP范围
	ipRange := strings.Split(targetIP, "/")
	if len(ipRange) != 2 {
		return nil, fmt.Errorf("IP范围格式不正确，应该是 x.x.x.x/x")
	}

	ip := ipRange[0]
	mask, err := strconv.Atoi(ipRange[1])
	if err != nil || mask < 0 || mask > 32 {
		return nil, fmt.Errorf("子网掩码无效")
	}

	// 计算起始和结束 IP
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return nil, fmt.Errorf("无效的 IP 地址")
	}

	ip4 := ipAddr.To4()
	if ip4 == nil {
		return nil, fmt.Errorf("只支持 IPv4 地址")
	}

	ipStart := ip4.To4()
	ipEnd := make(net.IP, len(ip4))
	copy(ipEnd, ip4)
	ipEnd[3] += 255

	// 启动扫描协程
	wg.Add(1)
	go func() {
		defer wg.Done()
		for ip := ipStart; ip[3] <= ipEnd[3]; ip[3]++ {
			ipStr := ip.String()
			if _, _, err := Ping(ipStr); err == nil {
				resultChan <- ipStr
			}
		}
		close(resultChan)
	}()

	// 等待所有协程完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []string
	for ip := range resultChan {
		results = append(results, ip)
	}

	return results, nil
}

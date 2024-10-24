package pkg

import (
	"fmt"
	"net"
	"sync"
)

// ScanPort 扫描端口
func ScanPort(target string, port int, resultChan chan int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", target, port))
	if err == nil {
		defer conn.Close()
		resultChan <- port // 发送端口号到结果管道
	}
}

// PortScan 扫描指定端口
func PortScan(target string, port int) {
	// 创建结果管道
	resultChan := make(chan int)

	// 启动多个协程扫描端口
	var wg sync.WaitGroup
	for i := 1; i <= port; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			ScanPort(target, port, resultChan)
		}(i)
	}

	// 等待所有协程完成
	wg.Wait()

	fmt.Println("扫描完成！")
	fmt.Println("正常退出")
}

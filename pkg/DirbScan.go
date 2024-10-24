package pkg

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
)

// Result 用于存储扫描结果
type Result struct {
	Path       string
	StatusCode int
}

// DicGet 读取字典
func DicGet(dicname string) ([]string, error) {
	file, err := os.Open(dicname)
	if err != nil {
		return nil, fmt.Errorf("无法打开字典文件: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取字典文件出错: %v", err)
	}

	fmt.Println("目录遍历中....")
	return lines, nil
}

// HttpCheck 判断web路径是否存在
func HttpCheck(url string, dics []string, resultsChan chan<- Result) {
	for _, dic := range dics {
		dicPath := url + dic
		resp, err := http.Get(dicPath)
		if err != nil {
			fmt.Printf("无法访问路径 %s: %v\n", dicPath, err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			resultsChan <- Result{Path: dicPath, StatusCode: resp.StatusCode}
		}
	}
	close(resultsChan) // 关闭通道，表示没有更多的数据发送
}

// DirbScan 执行目录扫描
func DirbScan(url string, dicname string) ([]Result, error) {
	if url == "" || dicname == "" {
		return nil, fmt.Errorf("-h 帮助")
	}

	// 检查字典文件是否存在
	dics, err := DicGet(dicname)
	if err != nil {
		return nil, err
	}

	resultsChan := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		HttpCheck(url, dics, resultsChan)
	}()

	var results []Result
	go func() {
		wg.Wait() // 等待 HttpCheck 完成
		for result := range resultsChan {
			results = append(results, result)
		}
	}()

	// 等待所有协程完成
	wg.Wait()

	return results, nil
}

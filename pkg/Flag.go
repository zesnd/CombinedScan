package pkg

import (
	"CombinedScan/utils"
	"flag"
	"fmt"
)

func Banner() {
	banner :=
		"      __                                       \n" +
			"     |  |           ______ ____ _____    ____  \n" +
			"     |  |  ______  /  ___// ___\\\\__  \\  /    \\ \n" +
			"     |  | /_____/  \\___ \\\\  \\___ / __ \\|   |  \\\n" +
			" /\\__|  |         /____  >\\___  >____  /___|  /\n" +
			" \\______|              \\/     \\/     \\/     \\/ "
	print(banner)
}

// CommandAction 定义命令行动作类型
type CommandAction func()

// SetFlags 设置命令行参数并返回命令行动作
func SetFlags(args []string) CommandAction {
	// 定义flagSet
	flagSet := flag.NewFlagSet("", flag.ExitOnError)

	fmt.Println("")

	// 定义参数字段
	var targetIP string       // 目标IP/C段
	var port int              // 目标端口
	var url string            // 目标URL
	var dicname string        // 字典名称目录
	var outputFilename string // 输出文件名
	var pocFile string        // POC文件名
	var dicFile string        // 批量验证文件名

	// 定义命令行参数
	flagSet.StringVar(&targetIP, "i", "", "目标IP/C段")
	flagSet.IntVar(&port, "p", 0, "扫描目标端口")
	flagSet.StringVar(&url, "u", "", "目标URL网址")
	flagSet.StringVar(&dicname, "d", "./dir/php.txt", "字典名称路径")
	flagSet.StringVar(&outputFilename, "o", "", "输出文件名称")
	flagSet.StringVar(&pocFile, "poc", "", "POC文件名")
	flagSet.StringVar(&dicFile, "dfile", "", "批量验证文件名")

	// 解析命令行参数
	err := flagSet.Parse(args[1:]) // 略过程序名
	if err != nil {
		fmt.Println("解析命令行参数失败:", err)
		flagSet.PrintDefaults()
		return nil
	}

	// 获取本地 IP 地址
	localIP, err := utils.GetLocalIP()
	if err != nil {
		fmt.Println("获取本地 IP 地址失败:", err)
		flagSet.PrintDefaults()
		return nil
	}

	// 根据参数判断执行的动作
	if targetIP != "" && port != 0 {
		return func() {
			// 执行端口扫描
			PortScan(targetIP, port)
		}
	}

	if targetIP != "" {
		return func() {
			// 执行 IP 扫描
			results, err := IPScan(targetIP)
			if err != nil {
				fmt.Println("IP扫描出错:", err)
				return
			}
			if outputFilename != "" {
				if err := utils.WriteResults(results, outputFilename); err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				for _, result := range results {
					fmt.Println(result)
				}
			}
		}
	}

	if url != "" && dicname != "" {
		return func() {
			// 执行目录扫描
			url, _ = utils.UrlHandler(url)
			results, _ := DirbScan(url, dicname)

			// Convert results to a slice of strings
			var resultStrings []string
			for _, result := range results {
				resultStrings = append(resultStrings, fmt.Sprintf("%s %d", result.Path, result.StatusCode))
			}

			// Write results to file or print to console
			if outputFilename != "" {
				if err := utils.WriteResults(resultStrings, outputFilename); err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				for _, result := range resultStrings {
					fmt.Println(result)
				}
			}
		}
	}

	if url != "" && pocFile != "" {
		return func() {
			url, _ = utils.UrlHandler(url)
			// 执行漏洞检测
			resultMessage, err := utils.RunPoc(pocFile, url, localIP)
			if err != nil {
				fmt.Println("漏洞检测出错:", err)
				return
			}
			if outputFilename != "" {
				if err := utils.WriteResults([]string{resultMessage}, outputFilename); err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println(resultMessage)
			}
		}
	}

	if dicFile != "" {
		return func() {
			// 读取文件中的 URL 列表
			urls, err := utils.ReadLines(dicFile)
			if err != nil {
				fmt.Println("读取文件出错:", err)
				return
			}

			var results []string
			for _, url := range urls {
				url, _ = utils.UrlHandler(url)
				// 执行漏洞检测
				resultMessage, err := utils.RunPoc(pocFile, url, localIP)
				if err != nil {
					fmt.Printf("URL %s 验证出错: %v\n", url, err)
					continue
				}
				results = append(results, resultMessage)
			}

			// 根据输出文件名决定输出位置
			if outputFilename != "" {
				if err := utils.WriteResults(results, outputFilename); err != nil {
					fmt.Println("写入文件出错:", err)
				}
			} else {
				for _, result := range results {
					fmt.Println(result)
				}
			}
		}
	}

	// 如果没有匹配到任何条件，打印帮助信息
	return func() {
		flagSet.PrintDefaults()
	}
}

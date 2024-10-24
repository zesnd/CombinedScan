// main.go

package main

import (
	"CombinedScan/pkg"
	"fmt"
	"os"
)

func main() {
	pkg.Banner()
	// 获取用户输入的命令行参数
	args := os.Args

	// 根据用户输入的参数设置命令行参数并获取相应的动作函数
	action := pkg.SetFlags(args)

	// 执行动作函数
	if action != nil {
		action()
	} else {
		fmt.Println("Invalid arguments")
	}
}

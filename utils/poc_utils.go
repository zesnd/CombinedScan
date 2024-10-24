package utils

import (
	"CombinedScan/poc"
	"fmt"
)

// RunPoc 运行指定的 POC
func RunPoc(pocFile, url, localIP string) (string, error) {
	pocFunctions := map[string]func(string, string) bool{
		"CVE-2021-44228": poc.CVE_2021_44228,
	}

	if pocFunc, exists := pocFunctions[pocFile]; exists {
		result := pocFunc(url, localIP)
		resultMessage := fmt.Sprintf("URL: %s, CVE: %s, Result: %v", url, pocFile, result)

		if result {
			fmt.Println("存在漏洞")
		} else {
			fmt.Println("未找到漏洞")
		}

		return resultMessage, nil
	} else {
		return "", fmt.Errorf("未知的POC类型: %s", pocFile)
	}
}

// utils/http_utils.go

package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// UrlHandler url处理函数
func UrlHandler(URL string) (string, error) {
	// 没有http前缀则添加
	if !strings.HasPrefix(URL, "http") {
		URL = "http://" + URL
	}

	// 删除路径并去掉结尾/
	targetURL, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("url解析失败: %w", err)
	}
	targetURL.Path = ""
	URL = strings.TrimSuffix(targetURL.String(), "/")

	return URL, nil
}

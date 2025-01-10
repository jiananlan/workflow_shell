package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func removeQuotes(str string) string {
	if len(str) > 1 && str[0] == '"' && str[len(str)-1] == '"' {
		return str[1 : len(str)-1]
	}
	return str
}

type ApiResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func checkPath(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			execPath, err := os.Executable()
			if err != nil {
			}
			path = filepath.Dir(execPath)
			return path
		}
	}
	if info.IsDir() {
		return path
	}
	path = filepath.Dir(path)
	return path
}

func req() {
	bff.WriteString("hhh\n")
	api := "https://api.github.com/repos/jiananlan/test/contents/test.txt?ref=main"
	if !stop {
		resp, err := http.Get(api)
		if err != nil {
			fmt.Println("请求失败:", err)
			return
		}
		defer resp.Body.Close()

		//fmt.Printf("HTTP 请求成功: 状态码 %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("解析 JSON 失败:", err)
			return
		}
		if response["content"] != nil {
			fmt.Println("failed to receive\n")
			return
		}
		encodedContent := response["content"].(string)

		decodedContent, err := base64.StdEncoding.DecodeString(encodedContent)
		if err != nil {
			fmt.Println("解码 base64 失败:", err)
			return
		}

		result := string(decodedContent)
		if result != last {
			muu.Lock()
			bff.WriteString(result + "\n")
			last = result
			muu.Unlock()
		}
	}
}

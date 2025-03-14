package main

import (
	"Jarvis/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// 提取 CMD 命令的函数
func extractCmdCommands(response string) []string {
	cmdRegex := regexp.MustCompile("(?i)```cmd\\s*([\\s\\S]+?)```")

	var commands []string
	matches := cmdRegex.FindAllStringSubmatch(response, -1)
	for _, match := range matches {
		if match[1] != "" {
			lines := strings.Split(match[1], "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed != "" && !strings.HasPrefix(trimmed, "//") {
					commands = append(commands, trimmed)
				}
			}
		}
	}
	return commands
}

// 处理来自 Python 的 HTTP 请求
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 解析 JSON 请求体
	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "无效的请求", http.StatusBadRequest)
		return
	}

	// 获取用户输入的消息
	message := requestData["message"]
	fmt.Println("接收到的消息:", message)

	// 读取 Tool.txt 作为系统提示词
	data, err := os.ReadFile("prompt/Tool.txt")
	if err != nil {
		log.Fatalf("Failed to read Tool.txt: %v", err) // 错误输出并终止程序
	}
	systemPrompt := string(data)

	// 创建 ChatUtils 实例
	chatUtils := utils.NewChatUtils(systemPrompt)

	// 发送请求
	response, err := chatUtils.SendChatRequest(message)
	if err != nil {
		log.Fatalf("Request failed: %v", err) // 错误输出并终止程序
	}

	// 输出完整响应
	fmt.Println("Response:", response)

	// 提取 CMD 命令
	cmdList := extractCmdCommands(response)

	// 输出提取的命令
	if len(cmdList) == 0 {
		fmt.Println("\n未找到 CMD 命令")
		w.Write([]byte("未找到 CMD 命令"))
	} else {
		fmt.Println("\n提取的 CMD 命令:")
		for _, cmd := range cmdList {
			fmt.Println(cmd)
		}

		// 执行 CMD 命令
		fmt.Println("\n正在执行命令...")
		err := utils.ExecuteCommands(cmdList)
		if err != nil {
			log.Fatalf("执行命令时出错: %v", err) // 错误输出并终止程序
		}
		fmt.Println("所有命令执行完成！")
		w.Write([]byte("命令执行完成"))
	}
}

func main() {
	// 设置 HTTP 路由
	http.HandleFunc("/execute-command", handleRequest)

	// 启动 HTTP 服务
	fmt.Println("Go 服务启动中...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("服务启动失败: %v", err) // 错误输出并终止程序
	}
}

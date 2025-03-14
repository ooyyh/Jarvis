package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 定义请求的结构体
const apiURL = "https://api.chatanywhere.tech/v1/chat/completions"
const apiKey = "" // 替换为你的实际 API 密钥

// RequestBody 代表请求体的结构
type RequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

// Message 代表单个消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatUtils 封装了与 Chat API 交互的工具
type ChatUtils struct {
	Messages     []Message // 存储多轮对话的消息
	SystemPrompt string    // 系统提示词
}

// NewChatUtils 创建 ChatUtils 实例（可选设置系统提示词）
func NewChatUtils(systemPrompt string) *ChatUtils {
	chatUtils := &ChatUtils{}

	// 如果提供了 systemPrompt，则添加到消息列表的最前面
	if systemPrompt != "" {
		chatUtils.Messages = append(chatUtils.Messages, Message{
			Role:    "system",
			Content: systemPrompt,
		})
		chatUtils.SystemPrompt = systemPrompt
	}

	return chatUtils
}

// SendChatRequest 发送请求到 Chat API，并返回响应
func (c *ChatUtils) SendChatRequest(message string) (string, error) {
	// 将用户消息添加到对话历史中
	c.Messages = append(c.Messages, Message{
		Role:    "user",
		Content: message,
	})

	// 创建请求体
	requestBody := RequestBody{
		Model:       "gpt-3.5-turbo",
		Messages:    c.Messages,
		Temperature: 0.7,
	}

	// 将请求体转换为 JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// 检查响应状态是否为 200 OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// 解析响应 JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// 提取模型回答
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("unexpected response format")
	}

	responseMessage, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	// 将模型的回答添加到对话历史中
	c.Messages = append(c.Messages, Message{
		Role:    "assistant",
		Content: responseMessage,
	})

	return responseMessage, nil
}

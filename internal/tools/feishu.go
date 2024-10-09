package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FeishuRequestBody 定义了飞书机器人的请求体结构
type FeishuRequestBody struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

// SendFeishuWebhook 发送 POST 请求到飞书机器人的 webhook URL
func SendFeishuWebhook(webhookURL, text string) error {
	// 创建请求体
	requestBody := FeishuRequestBody{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: text,
		},
	}

	// 将请求体编码为 JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// 发送 POST 请求
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send webhook: %s", body)
	}

	return nil
}

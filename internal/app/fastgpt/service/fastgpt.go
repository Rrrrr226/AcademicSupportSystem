package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FastGPTClient FastGPT 客户端
type FastGPTClient struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// NewFastGPTClient 创建 FastGPT 客户端
func NewFastGPTClient(baseURL, apiKey string) *FastGPTClient {
	return &FastGPTClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ForwardRequest 转发请求到 FastGPT
func (c *FastGPTClient) ForwardRequest(method, path string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// 发送请求
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// ForwardRequestWithQuery 带查询参数转发请求到 FastGPT
func (c *FastGPTClient) ForwardRequestWithQuery(method, path string, queryParams map[string]string) ([]byte, int, error) {
	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	// 添加查询参数
	if len(queryParams) > 0 {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// 发送请求
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// ForwardStreamRequest 转发流式请求到 FastGPT，返回响应对象用于流式读取
func (c *FastGPTClient) ForwardStreamRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	// 发送请求，使用无超时的客户端
	streamClient := &http.Client{
		Timeout: 0, // 流式请求不设置超时
	}
	resp, err := streamClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	return resp, nil
}

// StreamReader 流式读取器
type StreamReader struct {
	scanner *bufio.Scanner
	resp    *http.Response
}

// NewStreamReader 创建流式读取器
func NewStreamReader(resp *http.Response) *StreamReader {
	return &StreamReader{
		scanner: bufio.NewScanner(resp.Body),
		resp:    resp,
	}
}

// Read 读取下一个数据块
func (sr *StreamReader) Read() (string, bool) {
	if sr.scanner.Scan() {
		return sr.scanner.Text(), true
	}
	return "", false
}

// Close 关闭流
func (sr *StreamReader) Close() error {
	return sr.resp.Body.Close()
}

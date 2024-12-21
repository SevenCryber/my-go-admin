package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SevenCryber/my-go-admin/config/task"
	"io"
	"net/http"
	"net/url"
)

// ResponseInfo 定义返回的响应信息结构体
type ResponseInfo struct {
	Body    string
	TraceID string
}

func Send(config task.APIConfig, param interface{}, contentType string) (*ResponseInfo, error) {
	var bodyReader *bytes.Reader
	var req *http.Request
	var err error

	// 根据请求类型构建请求体或查询参数
	switch config.Method {
	case http.MethodGet:
		// 对于GET请求，参数应作为查询参数添加到URL
		u, err := url.Parse(config.Url)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %v", err)
		}

		q := u.Query()
		switch p := param.(type) {
		case string:
			// 如果param是字符串，则直接解析为查询参数
			extraQ, err := url.ParseQuery(p)
			if err != nil {
				return nil, fmt.Errorf("failed to parse query string: %v", err)
			}
			for k, v := range extraQ {
				q[k] = v
			}
		case map[string]interface{}:
			// 如果param是map，则将键值对添加到查询参数中
			for k, v := range p {
				q.Set(k, fmt.Sprint(v))
			}
		}
		u.RawQuery = q.Encode()
		req, err = http.NewRequest(config.Method, u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
	default:
		// 根据 contentType 决定如何处理请求体
		if contentType == "" || contentType == "json" {
			// 对于其他请求方法，默认尝试JSON序列化
			if param != nil {
				jsonData, err := json.Marshal(param)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal JSON: %v", err)
				}
				bodyReader = bytes.NewReader(jsonData)
			}
			req, err = http.NewRequest(config.Method, config.Url, bodyReader)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
		} else if contentType == "x-www-form-urlencoded" {
			// 对于表单数据，使用url.Values构建请求体或直接使用已编码的字符串
			switch p := param.(type) {
			case string:
				// 如果param已经是编码好的字符串，直接使用
				bodyReader = bytes.NewReader([]byte(p))
			case url.Values:
				// 如果param是url.Values，则编码后使用
				bodyReader = bytes.NewReader([]byte(p.Encode()))
			case map[string]interface{}:
				// 如果param是map，则转换为url.Values并编码
				values := url.Values{}
				for k, v := range p {
					values.Add(k, fmt.Sprint(v))
				}
				bodyReader = bytes.NewReader([]byte(values.Encode()))
			default:
				return nil, fmt.Errorf("unsupported param type for form data")
			}
			req, err = http.NewRequest(config.Method, config.Url, bodyReader)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			return nil, fmt.Errorf("unsupported content type: %s", contentType)
		}
		req.Header.Set("role-id", "admin")
		req.Header.Set("group-id", "123")

	}

	req.Header.Set("Cookie", "heat=123456")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 提取traceId
	traceID := resp.Header.Get("traceId")

	// 读取响应体内容，并加入缓冲区
	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 返回响应信息
	return &ResponseInfo{
		Body:    buffer.String(),
		TraceID: traceID,
	}, nil
}

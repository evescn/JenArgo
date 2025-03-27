package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
)

var Request request

type request struct {
}

// HttpRequest 统一的 HTTP 请求方法
func (*request) HttpRequest(method, service, url, privateToken string, body interface{}) ([]byte, error) {
	client := &http.Client{}

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			zap.L().Error("JSON 序列化失败: " + err.Error())
			return nil, err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		zap.L().Error("HTTP 请求创建失败: " + err.Error())
		return nil, err
	}

	// 设置请求头
	switch service {
	case "gitlab":
		req.Header.Set("PRIVATE-TOKEN", privateToken)
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	case "jenkins":
		// 兼容 jenkins post 请求
		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(privateToken))
		req.Header.Add("Authorization", authHeader)
	case "argo":
		// 兼容 jenkins post 请求
		//argoCDHeader := "Bearer " + base64.StdEncoding.EncodeToString([]byte(privateToken))
		//req.Header.Add("Authorization", argoCDHeader)
		req.Header.Set("Authorization", "Bearer "+privateToken)
		//argoCDHeader2 := "argocd.token=" + base64.StdEncoding.EncodeToString([]byte(privateToken))
		//req.Header.Add("Cookie", argoCDHeader2)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("HTTP 请求失败: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("读取响应失败: " + err.Error())
		return nil, err
	}

	return bodyBytes, nil
}

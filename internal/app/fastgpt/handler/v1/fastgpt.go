package v1

import (
	"HelpStudent/config"
	"HelpStudent/core/auth"
	"HelpStudent/core/logx"
	"HelpStudent/core/middleware/response"
	"HelpStudent/internal/app/fastgpt/dto"
	"HelpStudent/internal/app/fastgpt/service"
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

// getFastGPTClient 获取 FastGPT 客户端
func getFastGPTClient() *service.FastGPTClient {
	cfg := config.GetConfig()
	return service.NewFastGPTClient(cfg.FastGPT.BaseURL, "")
}

// HandleChatCompletion 处理聊天补全请求
func HandleChatCompletion(c flamego.Context, r flamego.Render, req dto.ChatCompletionRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	// 如果是流式请求
	if req.Stream {
		handleStreamChat(c, req)
		return
	}

	// 非流式请求
	respBody, statusCode, err := getFastGPTClient().ForwardRequest("POST", "/v1/chat/completions", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	// 直接返回 FastGPT 的响应
	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// handleStreamChat 处理流式聊天
func handleStreamChat(c flamego.Context, req dto.ChatCompletionRequest) {
	resp, err := getFastGPTClient().ForwardStreamRequest("POST", "/v1/chat/completions", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		c.ResponseWriter().Header().Set("Content-Type", "text/event-stream")
		c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.ResponseWriter(), "event: error\ndata: {\"error\":\"%s\"}\n\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d", resp.StatusCode)
		c.ResponseWriter().Header().Set("Content-Type", "text/event-stream")
		c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.ResponseWriter(), "event: error\ndata: {\"error\":\"FastGPT API 调用失败\"}\n\n")
		return
	}

	// 设置 SSE 响应头
	w := c.ResponseWriter()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	// 确保支持 Flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		logx.SystemLogger.CtxError(c.Request().Context(), "ResponseWriter does not support Flusher")
		return
	}

	// 读取并转发流式响应
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// 跳过空行
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// 直接转发 SSE 格式的数据
		fmt.Fprintf(w, "%s\n", line)
		flusher.Flush()

		// 检查是否是结束标记
		if strings.Contains(line, "[DONE]") {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), "Stream read error", err)
	}
}

// HandleGetHistories 获取聊天历史列表
func HandleGetHistories(c flamego.Context, r flamego.Render, req dto.GetHistoriesRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequest("POST", "/core/chat/history/getHistories", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// HandleUpdateHistory 更新聊天会话
func HandleUpdateHistory(c flamego.Context, r flamego.Render, req dto.UpdateHistoryRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequest("POST", "/core/chat/history/updateHistory", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// HandleDelHistory 删除聊天会话
func HandleDelHistory(c flamego.Context, r flamego.Render, authInfo auth.Info) {
	appId := c.Query("appId")
	chatId := c.Query("chatId")

	if appId == "" || chatId == "" {
		response.HTTPFail(r, 400001, "缺少必要参数")
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequestWithQuery("DELETE", "/core/chat/history/delHistory", map[string]string{
		"appId":  appId,
		"chatId": chatId,
	})
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// HandleGetPaginationRecords 获取聊天记录
func HandleGetPaginationRecords(c flamego.Context, r flamego.Render, req dto.GetPaginationRecordsRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequest("POST", "/core/chat/getPaginationRecords", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// HandleCreateDataset 创建数据集
func HandleCreateDataset(c flamego.Context, r flamego.Render, req dto.DatasetCreateRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequest("POST", "/core/dataset/create", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// HandleListDatasets 列出数据集
func HandleListDatasets(c flamego.Context, r flamego.Render, req dto.DatasetListRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequest("POST", "/core/dataset/list", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// HandleGetDatasetDetail 获取数据集详情
func HandleGetDatasetDetail(c flamego.Context, r flamego.Render, authInfo auth.Info) {
	id := c.Query("id")
	if id == "" {
		response.HTTPFail(r, 400001, "缺少必要参数")
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequestWithQuery("GET", "/core/dataset/detail", map[string]string{
		"id": id,
	})
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// HandleDeleteDataset 删除数据集
func HandleDeleteDataset(c flamego.Context, r flamego.Render, authInfo auth.Info) {
	id := c.Query("id")
	if id == "" {
		response.HTTPFail(r, 400001, "缺少必要参数")
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequestWithQuery("DELETE", "/core/dataset/delete", map[string]string{
		"id": id,
	})
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// HandleCreateCollectionText 从文本创建集合
func HandleCreateCollectionText(c flamego.Context, r flamego.Render, req dto.CreateCollectionTextRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequest("POST", "/core/dataset/collection/create/text", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// HandleCreateCollectionLink 从链接创建集合
func HandleCreateCollectionLink(c flamego.Context, r flamego.Render, req dto.CreateCollectionLinkRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequest("POST", "/core/dataset/collection/create/link", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// HandlePushData 推送数据到集合
func HandlePushData(c flamego.Context, r flamego.Render, req dto.PushDataRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequest("POST", "/core/dataset/data/pushData", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

// HandleSearchTest 搜索测试
func HandleSearchTest(c flamego.Context, r flamego.Render, req dto.SearchTestRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	respBody, statusCode, err := getFastGPTClient().ForwardRequest("POST", "/core/dataset/searchTest", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if statusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d, body=%s", statusCode, string(respBody))
		response.HTTPFail(r, 500001, "FastGPT API 调用失败")
		return
	}

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(respBody)
}

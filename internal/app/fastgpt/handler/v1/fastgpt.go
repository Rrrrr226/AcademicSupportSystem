package v1

import (
	"HelpStudent/config"
	"HelpStudent/core/auth"
	"HelpStudent/core/logx"
	"HelpStudent/core/middleware/response"
	"HelpStudent/internal/app/fastgpt/dao"
	"HelpStudent/internal/app/fastgpt/dto"
	"HelpStudent/internal/app/fastgpt/service"
	"bufio"
	"net/http"
	"strings"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

// getFastGPTClient 获取 FastGPT 客户端（使用指定的 API Key）
func getFastGPTClient(apiKey string) *service.FastGPTClient {
	cfg := config.GetConfig()
	return service.NewFastGPTClient(cfg.FastGPT.BaseURL, apiKey)
}

// HandleChatCompletion 处理聊天补全请求
func HandleChatCompletion(c flamego.Context, r flamego.Render, req dto.ChatCompletionRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}
	// TODO检查这个用户是否可以使用这个app

	// 根据 fastgptAppId 获取对应的 API Key
	app, err := dao.FastgptApp.GetAppByID(req.FastgptAppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.HTTPFail(r, 400013, "应用不存在或已禁用")
		return
	}

	// 如果是流式请求，返回提示使用流式接口
	if req.Stream {
		response.HTTPFail(r, 400014, "请使用流式接口 /v1/chat/completions/stream")
		return
	}

	// 非流式请求
	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequest("POST", "/v1/chat/completions", req)
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

// HandleStreamChatCompletion 处理流式聊天补全请求（使用 flamego/sse）
func HandleStreamChatCompletion(c flamego.Context, req dto.ChatCompletionRequest, errs binding.Errors, authInfo auth.Info, msg chan<- *dto.SSEMessage) {
	if errs != nil {
		msg <- &dto.SSEMessage{Data: `{"error":"参数错误"}`, Event: "error"}
		return
	}

	// 根据 fastgptAppId 获取对应的 API Key
	app, err := dao.FastgptApp.GetAppByID(req.FastgptAppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		msg <- &dto.SSEMessage{Data: `{"error":"应用不存在或已禁用"}`, Event: "error"}
		return
	}

	// 强制设置为流式模式
	req.Stream = true

	// 发起流式请求
	resp, err := getFastGPTClient(app.APIKey).ForwardStreamRequest("POST", "/v1/chat/completions", req)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		msg <- &dto.SSEMessage{Data: `{"error":"请求失败"}`, Event: "error"}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logx.SystemLogger.CtxError(c.Request().Context(), "FastGPT API error: status=%d", resp.StatusCode)
		msg <- &dto.SSEMessage{Data: `{"error":"FastGPT API 调用失败"}`, Event: "error"}
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

		// 解析 SSE 格式数据
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			msg <- &dto.SSEMessage{Data: data}

			// 检查是否是结束标记
			if data == "[DONE]" {
				break
			}
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

	// 使用 appId 获取 API Key（FastGPT 也需要 appId 字段）
	app, err := dao.FastgptApp.GetAppByID(req.AppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequest("POST", "/core/chat/history/getHistories", req)
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

	// 使用 appId 获取 API Key（FastGPT 也需要 appId 字段）
	app, err := dao.FastgptApp.GetAppByID(req.AppId)

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequest("POST", "/core/chat/history/updateHistory", req)
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

	app, err := dao.FastgptApp.GetAppByID(appId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequestWithQuery("DELETE", "/core/chat/history/delHistory", map[string]string{
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

	// 使用 appId 获取 API Key（FastGPT 也需要 appId 字段）
	app, err := dao.FastgptApp.GetAppByID(req.AppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequest("POST", "/core/chat/getPaginationRecords", req)
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

	app, err := dao.FastgptApp.GetAppByID(req.FastgptAppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequest("POST", "/core/dataset/create", req)
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

	app, err := dao.FastgptApp.GetAppByID(req.FastgptAppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequest("POST", "/core/dataset/list", req)
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
	fastgptAppId := c.Query("fastgptAppId")
	if id == "" || fastgptAppId == "" {
		response.HTTPFail(r, 400001, "缺少必要参数")
		return
	}

	app, err := dao.FastgptApp.GetAppByID(fastgptAppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequestWithQuery("GET", "/core/dataset/detail", map[string]string{
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
	fastgptAppId := c.Query("fastgptAppId")
	if id == "" || fastgptAppId == "" {
		response.HTTPFail(r, 400001, "缺少必要参数")
		return
	}

	app, err := dao.FastgptApp.GetAppByID(fastgptAppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequestWithQuery("DELETE", "/core/dataset/delete", map[string]string{
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

	app, err := dao.FastgptApp.GetAppByID(req.FastgptAppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequest("POST", "/core/dataset/collection/create/text", req)
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

	app, err := dao.FastgptApp.GetAppByID(req.FastgptAppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequest("POST", "/core/dataset/collection/create/link", req)
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

	app, err := dao.FastgptApp.GetAppByID(req.FastgptAppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequest("POST", "/core/dataset/data/pushData", req)
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

	app, err := dao.FastgptApp.GetAppByID(req.FastgptAppId)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	respBody, statusCode, err := getFastGPTClient(app.APIKey).ForwardRequest("POST", "/core/dataset/searchTest", req)
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

package router

import (
	"HelpStudent/core/middleware/web"
	"HelpStudent/internal/app/fastgpt/dto"
	handler "HelpStudent/internal/app/fastgpt/handler/v1"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"github.com/flamego/sse"
)

func AppFastgptInit(e *flamego.Flame) {
	// FastGPT API 转发接口（需要登录）
	e.Group("/fastgpt", func() {
		// Chat 接口 - 非流式
		e.Post("/v1/chat/completions", binding.JSON(dto.ChatCompletionRequest{}), handler.HandleChatCompletion)
		// Chat 接口 - 流式输出（使用 flamego/sse）
		e.Post("/v1/chat/completions/stream", binding.JSON(dto.ChatCompletionRequest{}), sse.Bind(dto.SSEMessage{}), handler.HandleStreamChatCompletion)

		// Chat History 接口
		e.Group("/core/chat", func() {
			e.Post("/history/getHistories", binding.JSON(dto.GetHistoriesRequest{}), handler.HandleGetHistories)
			e.Post("/history/updateHistory", binding.JSON(dto.UpdateHistoryRequest{}), handler.HandleUpdateHistory)
			e.Delete("/history/delHistory", handler.HandleDelHistory)
			e.Post("/getPaginationRecords", binding.JSON(dto.GetPaginationRecordsRequest{}), handler.HandleGetPaginationRecords)
		})

		// Dataset 接口
		e.Group("/core/dataset", func() {
			e.Post("/create", binding.JSON(dto.DatasetCreateRequest{}), handler.HandleCreateDataset)
			e.Post("/list", binding.JSON(dto.DatasetListRequest{}), handler.HandleListDatasets)
			e.Get("/detail", handler.HandleGetDatasetDetail)
			e.Delete("/delete", handler.HandleDeleteDataset)

			// Collection 接口
			e.Group("/collection", func() {
				e.Post("/create/text", binding.JSON(dto.CreateCollectionTextRequest{}), handler.HandleCreateCollectionText)
				e.Post("/create/link", binding.JSON(dto.CreateCollectionLinkRequest{}), handler.HandleCreateCollectionLink)
			})

			// Data 接口
			e.Post("/data/pushData", binding.JSON(dto.PushDataRequest{}), handler.HandlePushData)
			e.Post("/searchTest", binding.JSON(dto.SearchTestRequest{}), handler.HandleSearchTest)
		})

		// App 管理接口
		e.Group("/apps", func() {
			e.Post("/create", binding.JSON(dto.CreateAppRequest{}), handler.HandleCreateApp)
			e.Post("/list", binding.JSON(dto.GetAppListRequest{}), handler.HandleGetAppList)
			e.Post("/update", binding.JSON(dto.UpdateAppRequest{}), handler.HandleUpdateApp)
			e.Post("/delete", binding.JSON(dto.DeleteAppRequest{}), handler.HandleDeleteApp)
		})
	}, web.Authorization)
}

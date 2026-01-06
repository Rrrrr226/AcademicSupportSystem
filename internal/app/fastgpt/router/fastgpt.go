package router

import (
	"HelpStudent/core/middleware/response"
	"HelpStudent/core/middleware/web"
	"HelpStudent/internal/app/fastgpt/dto"
	handler "HelpStudent/internal/app/fastgpt/handler/v1"
	"errors"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func AppFastgptInit(e *flamego.Flame) {
	e.Get("/fastgpt/v1", func(r flamego.Render) {
		response.HTTPSuccess(r, map[string]any{
			"message": "fastgpt Init Success",
		})
	})

	e.Get("/fastgpt/v1/err", func(r flamego.Render) {
		response.HTTPFail(r, 500000, "fastgpt Init test error", errors.New("this is err"))
	})

	// FastGPT API 转发接口（需要登录）
	e.Group("/fastgpt", func() {
		// Chat 接口 - 支持流式输出
		e.Post("/v1/chat/completions", binding.JSON(dto.ChatCompletionRequest{}), handler.HandleChatCompletion)

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

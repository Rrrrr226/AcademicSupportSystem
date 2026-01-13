package dto

// ChatCompletionRequest Chat 请求
type ChatCompletionRequest struct {
	FastgptAppId string                 `json:"fastgptAppId" binding:"Required"` // 用于获取 API Key，不转发给 FastGPT
	ChatId       string                 `json:"chatId"`
	Stream       bool                   `json:"stream"`
	Detail       bool                   `json:"detail"`
	Variables    map[string]interface{} `json:"variables"`
	Messages     []Message              `json:"messages" binding:"Required"`
	CustomUid    string                 `json:"customUid"`
	ShareId      string                 `json:"shareId"`
	OutLinkUid   string                 `json:"outLinkUid"`
}

type Message struct {
	Role    string      `json:"role" binding:"Required"`
	Content interface{} `json:"content" binding:"Required"`
}

// GetHistoriesRequest 获取聊天历史列表请求
// FastGPT API 需要 appId 字段，所以保留
type GetHistoriesRequest struct {
	FastgptAppId string `json:"fastgptAppId" binding:"Required"`
	Offset       int    `json:"offset"`
	PageSize     int    `json:"pageSize"`
	Source       string `json:"source"`
	ShareId      string `json:"shareId"`
	OutLinkUid   string `json:"outLinkUid"`
}

// UpdateHistoryRequest 更新聊天会话请求
// FastGPT API 需要 appId 字段，所以保留
type UpdateHistoryRequest struct {
	AppId       string `json:"appId" binding:"Required"`
	ChatId      string `json:"chatId" binding:"Required"`
	CustomTitle string `json:"customTitle"`
	Top         *bool  `json:"top"`
}

// GetPaginationRecordsRequest 获取聊天记录请求
// FastGPT API 需要 appId 字段，所以保留
type GetPaginationRecordsRequest struct {
	FastgptAppId        string `json:"fastgptAppId" binding:"Required"`
	AppId               string `json:"appId" binding:"Required"`
	ChatId              string `json:"chatId" binding:"Required"`
	Offset              int    `json:"offset"`
	PageSize            int    `json:"pageSize"`
	LoadCustomFeedbacks bool   `json:"loadCustomFeedbacks"`
}

// DatasetCreateRequest 创建数据集请求
// FastGPT API 不需要 appId，使用 fastgptAppId 获取 API Key
type DatasetCreateRequest struct {
	FastgptAppId string  `json:"fastgptAppId" binding:"Required"`
	ParentId     *string `json:"parentId"`
	Type         string  `json:"type"`
	Name         string  `json:"name" binding:"Required"`
	Intro        string  `json:"intro"`
	Avatar       string  `json:"avatar"`
	VectorModel  string  `json:"vectorModel"`
	AgentModel   string  `json:"agentModel"`
}

// DatasetListRequest 数据集列表请求
// FastGPT API 不需要 appId，使用 fastgptAppId 获取 API Key
type DatasetListRequest struct {
	FastgptAppId string  `json:"fastgptAppId" binding:"Required"`
	ParentId     *string `json:"parentId"`
}

// CreateCollectionTextRequest 从文本创建集合请求
// FastGPT API 不需要 appId，使用 fastgptAppId 获取 API Key
type CreateCollectionTextRequest struct {
	FastgptAppId     string `json:"fastgptAppId" binding:"Required"`
	Text             string `json:"text" binding:"Required"`
	DatasetId        string `json:"datasetId" binding:"Required"`
	Name             string `json:"name" binding:"Required"`
	TrainingType     string `json:"trainingType" binding:"Required"`
	ChunkSettingMode string `json:"chunkSettingMode"`
}

// CreateCollectionLinkRequest 从链接创建集合请求
// FastGPT API 不需要 appId，使用 fastgptAppId 获取 API Key
type CreateCollectionLinkRequest struct {
	FastgptAppId string                 `json:"fastgptAppId" binding:"Required"`
	Link         string                 `json:"link" binding:"Required"`
	DatasetId    string                 `json:"datasetId" binding:"Required"`
	TrainingType string                 `json:"trainingType" binding:"Required"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// PushDataRequest 推送数据请求
// FastGPT API 不需要 appId，使用 fastgptAppId 获取 API Key
type PushDataRequest struct {
	FastgptAppId string     `json:"fastgptAppId" binding:"Required"`
	CollectionId string     `json:"collectionId" binding:"Required"`
	TrainingType string     `json:"trainingType"`
	Data         []DataItem `json:"data" binding:"Required"`
}

type DataItem struct {
	Q       string                   `json:"q" binding:"Required"`
	A       string                   `json:"a"`
	Indexes []map[string]interface{} `json:"indexes"`
}

// SearchTestRequest 搜索测试请求
// FastGPT API 不需要 appId，使用 fastgptAppId 获取 API Key
type SearchTestRequest struct {
	FastgptAppId string  `json:"fastgptAppId" binding:"Required"`
	DatasetId    string  `json:"datasetId" binding:"Required"`
	Text         string  `json:"text" binding:"Required"`
	Limit        int     `json:"limit"`
	Similarity   float64 `json:"similarity"`
	SearchMode   string  `json:"searchMode"`
}

// === FastGPT App 管理相关 DTO ===

// CreateAppRequest 创建应用请求
type CreateAppRequest struct {
	AppName     string `json:"appName" binding:"Required"`
	AppId       string `json:"appId" binding:"Required"`
	ShareId     string `json:"shareId"`
	APIKey      string `json:"apiKey" binding:"Required"`
	Description string `json:"description"`
}

// UpdateAppRequest 更新应用请求
type UpdateAppRequest struct {
	ID          string `json:"id" binding:"Required"`
	AppName     string `json:"appName"`
	AppId       string `json:"appId"`
	ShareId     string `json:"shareId"`
	APIKey      string `json:"apiKey"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
}

// DeleteAppRequest 删除应用请求
type DeleteAppRequest struct {
	ID string `json:"id" binding:"Required"`
}

// GetAppListRequest 获取应用列表请求
type GetAppListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// AppItem 应用列表项
type AppItem struct {
	ID          string `json:"id"`
	AppName     string `json:"appName"`
	AppId       string `json:"appId"`
	ShareId     string `json:"shareId"`
	APIKey      string `json:"apiKey"`
	Description string `json:"description"`
	CreatedBy   string `json:"createdBy"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// AppListResponse 应用列表响应
type AppListResponse struct {
	Apps  []AppItem `json:"apps"`
	Total int64     `json:"total"`
}

// CreateAppResponse 创建应用响应
type CreateAppResponse struct {
	Name string `json:"name"`
}

// SSEMessage SSE 消息结构（用于流式输出）
type SSEMessage struct {
	Data  string `json:"data"`
	Event string `json:"event,omitempty"`
}

// GetCollectionQuoteRequest 获取集合引用请求
type GetCollectionQuoteRequest struct {
	FastgptAppId   string `json:"fastgptAppId" binding:"Required"` // 用于获取 API Key
	InitialId      string `json:"initialId"`
	InitialIndex   int    `json:"initialIndex"`
	PageSize       int    `json:"pageSize"`
	CollectionId   string `json:"collectionId" binding:"Required"`
	ChatItemDataId string `json:"chatItemDataId"`
	ChatId         string `json:"chatId"`
	AppId          string `json:"appId"`
	ShareId        string `json:"shareId"`
	OutLinkUid     string `json:"outLinkUid"`
}

package dto

// ChatCompletionRequest Chat 请求
type ChatCompletionRequest struct {
	ChatId    string                 `json:"chatId"`
	Stream    bool                   `json:"stream"`
	Detail    bool                   `json:"detail"`
	Variables map[string]interface{} `json:"variables"`
	Messages  []Message              `json:"messages" binding:"Required"`
}

type Message struct {
	Role    string      `json:"role" binding:"Required"`
	Content interface{} `json:"content" binding:"Required"`
}

// GetHistoriesRequest 获取聊天历史列表请求
type GetHistoriesRequest struct {
	AppId    string `json:"appId"`
	Offset   int    `json:"offset"`
	PageSize int    `json:"pageSize"`
	Source   string `json:"source"`
}

// UpdateHistoryRequest 更新聊天会话请求
type UpdateHistoryRequest struct {
	AppId       string `json:"appId" binding:"Required"`
	ChatId      string `json:"chatId" binding:"Required"`
	CustomTitle string `json:"customTitle"`
	Top         *bool  `json:"top"`
}

// GetPaginationRecordsRequest 获取聊天记录请求
type GetPaginationRecordsRequest struct {
	AppId               string `json:"appId" binding:"Required"`
	ChatId              string `json:"chatId" binding:"Required"`
	Offset              int    `json:"offset"`
	PageSize            int    `json:"pageSize"`
	LoadCustomFeedbacks bool   `json:"loadCustomFeedbacks"`
}

// DatasetCreateRequest 创建数据集请求
type DatasetCreateRequest struct {
	ParentId    *string `json:"parentId"`
	Type        string  `json:"type"`
	Name        string  `json:"name" binding:"Required"`
	Intro       string  `json:"intro"`
	Avatar      string  `json:"avatar"`
	VectorModel string  `json:"vectorModel"`
	AgentModel  string  `json:"agentModel"`
}

// DatasetListRequest 数据集列表请求
type DatasetListRequest struct {
	ParentId *string `json:"parentId"`
}

// CreateCollectionTextRequest 从文本创建集合请求
type CreateCollectionTextRequest struct {
	Text             string `json:"text" binding:"Required"`
	DatasetId        string `json:"datasetId" binding:"Required"`
	Name             string `json:"name" binding:"Required"`
	TrainingType     string `json:"trainingType" binding:"Required"`
	ChunkSettingMode string `json:"chunkSettingMode"`
}

// CreateCollectionLinkRequest 从链接创建集合请求
type CreateCollectionLinkRequest struct {
	Link         string                 `json:"link" binding:"Required"`
	DatasetId    string                 `json:"datasetId" binding:"Required"`
	TrainingType string                 `json:"trainingType" binding:"Required"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// PushDataRequest 推送数据请求
type PushDataRequest struct {
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
type SearchTestRequest struct {
	DatasetId  string  `json:"datasetId" binding:"Required"`
	Text       string  `json:"text" binding:"Required"`
	Limit      int     `json:"limit"`
	Similarity float64 `json:"similarity"`
	SearchMode string  `json:"searchMode"`
}

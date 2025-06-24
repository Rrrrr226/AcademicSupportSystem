package sls

import (
	"HelpStudent/core/logx"
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"time"
)

type Client struct {
	client sls.ClientInterface
}

type LogSearch struct {
	project  string
	logStore string
	topic    string
	from     int64
	to       int64
	query    string
	line     int64
	offset   int64
	reverse  bool
}

func NewClient(conf logx.SlsSinkConf) (*Client, error) {
	u, err := ParseURL(conf.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid aliyunsls-url: %w", err)
	}
	credentialsProvider := sls.NewStaticCredentialsProvider(u.AccessKeyID, u.AccessKeySecret, "")
	client := sls.CreateNormalInterfaceV2(u.Endpoint, credentialsProvider)
	return &Client{
		client: client,
	}, nil
}

func (c *Client) BuildLogSearch() *LogSearch {
	return &LogSearch{
		from: time.Now().Unix() - 3600,
		to:   time.Now().Unix(),
		// 其他的都用默认值
	}
}

func (c *Client) Search(search *LogSearch) (*sls.GetLogsResponse, error) {
	response, err := c.client.GetLogs(search.project, search.logStore, search.topic, search.from,
		search.to, search.query, search.line, search.offset, search.reverse)
	if err != nil {
		return nil, fmt.Errorf("get logs failed: %w", err)
	}
	return response, nil
}

func (c *Client) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

func (s *LogSearch) WithQuery(query string) *LogSearch {
	s.query = query
	return s
}
func (s *LogSearch) WithProject(project string) *LogSearch {
	s.project = project
	return s
}

func (s *LogSearch) WithLogStore(logStore string) *LogSearch {
	s.logStore = logStore
	return s
}

func (s *LogSearch) WithTopic(topic string) *LogSearch {
	s.topic = topic
	return s
}

func (s *LogSearch) WithSearchStartTime(from int64) *LogSearch {
	s.from = from
	return s
}

func (s *LogSearch) WithSearchEndTime(to int64) *LogSearch {
	s.to = to
	return s
}

func (s *LogSearch) WithLine(line int64) *LogSearch {
	s.line = line
	return s
}

func (s *LogSearch) WithOffset(offset int64) *LogSearch {
	s.offset = offset
	return s
}

func (s *LogSearch) WithReverse(reverse bool) *LogSearch {
	s.reverse = reverse
	return s
}

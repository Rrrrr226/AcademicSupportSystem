package fileServer

import (
	"HelpStudent/core/logx"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"io"
	"net/url"
	"path"
	"strings"
)

type AliOSS struct {
	client *oss.Client
	bucket *oss.Bucket
	config Config
}

func NewAliOSS(cfg Config) (FileClient, error) {
	var client AliOSS
	client.config = cfg
	if cfg.AccessKeyId == "" || cfg.AccessKeySecret == "" || cfg.EndPoint == "" || cfg.BucketName == "" {
		return nil, errors.New("阿里云OSS未配置完整，请检查配置文件")
	}
	cli, err := oss.New(cfg.EndPoint, cfg.AccessKeyId, cfg.AccessKeySecret)
	if err != nil {
		return nil, errors.Wrap(err, "无法创建OSS客户端")
	}
	client.client = cli
	bucket, err := cli.Bucket(cfg.BucketName)
	if err != nil {
		return nil, errors.Wrap(err, "无法获取OSS Bucket 存储空间")
	}
	client.bucket = bucket
	return &client, nil
}

func (ali *AliOSS) UploadFile(file []byte, fileName string) (string, error) {
	data := bytes.NewReader(file)
	err := ali.bucket.PutObject(path.Join(ali.config.Prefix, fileName), data)
	if err != nil {
		return "", errors.Wrap(err, "上传文件到OSS失败")
	}
	return ali.GetPreviewLink(fileName)
}

func (ali *AliOSS) UploadFileFromIO(fd io.Reader, fileName string) (string, error) {
	err := ali.bucket.PutObject(path.Join(ali.config.Prefix, fileName), fd)
	if err != nil {
		return "", errors.Wrap(err, "上传文件到OSS失败")
	}
	return ali.GetPreviewLink(fileName)
}

func (ali *AliOSS) ReadAll(fileName string) ([]byte, error) {
	file, err := ali.bucket.GetObject(path.Join(ali.config.Prefix, fileName))
	if err != nil {
		return nil, errors.Wrap(err, "读取文件失败")
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "读取文件失败")
	}
	return data, nil
}

func (ali *AliOSS) DeleteFile(fileName string) error {
	err := ali.bucket.DeleteObject(path.Join(ali.config.Prefix, fileName))
	if err != nil {
		return errors.Wrap(err, "删除文件失败")
	}
	return nil
}

func (ali *AliOSS) ReadDir(dir string) ([]DirItem, error) {
	continueToken := ""
	prefix := oss.Prefix(path.Join(ali.config.Prefix, dir))
	result := make([]DirItem, 0)
	for {
		lsRes, err := ali.bucket.ListObjectsV2(prefix, oss.ContinuationToken(continueToken))
		if err != nil {
			return nil, errors.Wrap(err, "读取文件夹失败")
		}
		for _, object := range lsRes.Objects {
			result = append(result, DirItem{
				Name:  object.Key,
				IsDir: false,
				Size:  object.Size,
			})
		}
		if lsRes.IsTruncated {
			prefix = oss.Prefix(lsRes.Prefix)
			continueToken = lsRes.NextContinuationToken
		} else {
			break
		}
	}
	return result, nil
}

func (ali *AliOSS) DownloadLink(fileName string) (string, error) {
	link, err := ali.bucket.SignURL(path.Join(ali.config.Prefix, fileName), oss.HTTPGet, 300)
	if err != nil {
		return "", errors.Wrap(err, "获取下载链接失败")
	}
	parse, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	parse.Scheme = ali.config.Schema
	parse.Host = ali.config.Host
	return parse.String(), nil
}

func (ali *AliOSS) UploadLink(fileName string, contentType string, callback CallbackConfig) (string, error) {
	callbackBody := "{\"contentMd5\": ${contentMd5},\"filename\": ${object},\"mimeType\":${mimeType},\"size\":${size}"

	callbackVar := make(map[string]string)
	for k, v := range callback {
		callbackBody += ",\"" + k + "\":${x:" + k + "}"
		callbackVar["x:"+k] = v
	}
	d, _ := json.Marshal(callbackVar)
	callbackBody += "}"

	link, err := ali.bucket.SignURL(path.Join(ali.config.Prefix, fileName), oss.HTTPPut, 300, oss.ContentType(contentType),
		oss.AddParam("callback", createCallbackString(&CallbackSign{
			CallbackUrl:      strings.Join(ali.config.CallbackUrls, ";"),
			CallbackBody:     callbackBody,
			CallbackBodyType: "application/json",
		})),
		oss.AddParam("callback-var", base64.StdEncoding.EncodeToString(d)),
	)
	if err != nil {
		return "", errors.Wrap(err, "获取上传链接失败")
	}
	parse, err := url.Parse(link)
	if err != nil {
		return link, err
	}
	parse.Scheme = ali.config.Schema
	parse.Host = ali.config.Host
	return parse.String(), nil
}

func (ali *AliOSS) StsToken(fileName string, action string) (string, error) {
	return "", errors.New("not implemented")
}

func (ali *AliOSS) GetPreviewLink(fileName string) (string, error) {
	if !ali.needConvert(fileName) {
		return ali.DownloadLink(fileName)
	}
	link, err := ali.bucket.SignURL(path.Join(ali.config.Prefix, fileName), oss.HTTPGet, 600, oss.Process("doc/preview,export_0,print_0"))
	if err != nil {
		return "", errors.Wrap(err, "获取预览链接失败")
	}
	parse, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	parse.Scheme = ali.config.Schema
	parse.Host = ali.config.Host
	return parse.String(), nil
}

func createCallbackString(data *CallbackSign) string {
	callbackStr, err := json.Marshal(data)
	if err != nil {
		logx.SystemLogger.Error("callback json make err:", err)
	}
	return base64.StdEncoding.EncodeToString(callbackStr)

}

type CallbackSign struct {
	CallbackUrl string `json:"callbackUrl"`
	//CallbackHost     string `json:"callbackHost"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

// https://help.aliyun.com/zh/oss/user-guide/overview-65
func (ali *AliOSS) needConvert(fileName string) bool {
	ext := path.Ext(fileName)
	switch ext {
	case ".doc", ".dot", ".wps", ".wpt", ".docx", ".dotx", ".docm", ".dotm", ".rtf", ".xls", ".xlt", ".et", ".xlsx", ".xltx", ".csv", ".xlsm", ".xltm", ".ppt", ".pptx", ".pptm", ".ppsx", ".ppsm", ".pps", ".potx", ".potm", ".dpt", ".dps", ".pdf":
		return true
	default:
		return false
	}
}

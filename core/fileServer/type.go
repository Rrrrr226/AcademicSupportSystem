package fileServer

import "io"

type DirItem struct {
	Name  string
	IsDir bool
	Size  int64
}

type CallbackConfig map[string]string

type FileClient interface {
	UploadFile(file []byte, fileName string) (string, error)
	UploadFileFromIO(fd io.Reader, fileName string) (string, error)
	ReadAll(fileName string) ([]byte, error)
	DeleteFile(fileName string) error
	ReadDir(dir string) ([]DirItem, error)
	DownloadLink(fileName string) (string, error)
	UploadLink(string, string, CallbackConfig) (string, error)
	StsToken(fileName string, action string) (string, error)
	//GeneratePreview(fileName string) error
	//HasPreview(fileName string) (bool, error)
	GetPreviewLink(fileName string) (string, error)
}

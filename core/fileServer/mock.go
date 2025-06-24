package fileServer

import (
	"github.com/pkg/errors"
	"io"
)

type Mock struct{}

func NewMock(Config) (FileClient, error) {
	return &Mock{}, nil
}

func (Mock) UploadFile([]byte, string) (string, error) {
	return "", errors.New("not implemented")
}

func (Mock) UploadFileFromIO(io.Reader, string) (string, error) {
	return "", errors.New("not implemented")
}

func (Mock) ReadAll(string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (Mock) DeleteFile(string) error {
	return errors.New("not implemented")
}

func (Mock) ReadDir(string) ([]DirItem, error) {
	return nil, errors.New("not implemented")
}

func (Mock) DownloadLink(string) (string, error) {
	return "", errors.New("not implemented")
}

func (Mock) UploadLink(string, string, CallbackConfig) (string, error) {
	return "", errors.New("not implemented")
}

func (Mock) StsToken(string, string) (string, error) {
	return "", errors.New("not implemented")
}

func (Mock) GeneratePreview(string) error {
	return errors.New("not implemented")
}

func (Mock) HasPreview(string) (bool, error) {
	return false, errors.New("not implemented")
}

func (Mock) GetPreviewLink(string) (string, error) {
	return "", errors.New("not implemented")
}

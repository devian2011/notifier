package storage

import (
	"errors"
	"net/url"
	"os"
	"strings"

	"notifications/internal/handler"
)

var (
	ErrMessageAlreadyExists = errors.New("message with same code already exists")
	ErrMessageNotExists     = errors.New("message with same code not exists")
)

type Config struct {
	Path string `json:"path" yaml:"path" `
}

func Factory(dsnStr string) (handler.TemplateStorage, error) {
	pwd, _ := os.Getwd()
	dsnStr = strings.ReplaceAll(dsnStr, "PWD", pwd)
	dsn, parseErr := url.Parse(dsnStr)
	if parseErr != nil {
		return nil, parseErr
	}

	switch dsn.Scheme {
	case "file":
		return NewStorage(dsn), nil
	default:
		return nil, errors.New("unknown storage")
	}
}

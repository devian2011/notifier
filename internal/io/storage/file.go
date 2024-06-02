package file

import (
	"errors"
	"fmt"
	"os"
	"strings"
	
	"github.com/kos-v/dsnparser"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"notifications/internal/dto"
	"notifications/internal/io/storage"
)

type Storage struct {
	dsn *dsnparser.DSN
}

func NewStorage(dsn *dsnparser.DSN) *Storage {
	return &Storage{dsn: dsn}
}

func (s *Storage) getFilePath(code string) string {
	code = strings.ReplaceAll(code, string(os.PathSeparator), "_")
	return fmt.Sprintf("%s/%s.yml", s.dsn.GetPath(), code)
}

func (s *Storage) List() ([]dto.MessageTmpl, error) {
	entries, entriesErr := os.ReadDir(s.dsn.GetPath())
	if entriesErr != nil {
		return nil, entriesErr
	}

	result := make([]dto.MessageTmpl, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fInfo, fErr := entry.Info()
		if fErr != nil {
			logrus.WithField("error", fErr).Errorf("error on file info")
			continue
		}
		filePath := fmt.Sprintf("%s/%s", s.dsn.GetPath(), fInfo.Name())

		msg := dto.MessageTmpl{}
		file, fileErr := os.ReadFile(filePath)
		if fileErr != nil {
			logrus.WithField("error", fileErr).Errorf("error on file read")
			continue
		}

		uErr := yaml.Unmarshal(file, &msg)

		if uErr != nil {
			logrus.WithField("error", uErr).Errorf("error on parse template")
			continue
		}

		result = append(result, msg)
	}

	return result, nil
}

func (s *Storage) Load(code string) (dto.MessageTmpl, error) {
	msg := dto.MessageTmpl{}
	file, fileErr := os.ReadFile(s.getFilePath(code))
	if fileErr != nil {
		return msg, fileErr
	}

	uErr := yaml.Unmarshal(file, &msg)

	return msg, uErr
}

func (s *Storage) Create(msg dto.MessageTmpl) error {
	filePath := s.getFilePath(msg.Code)
	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		return storage.ErrMessageAlreadyExists
	}
	file, fileErr := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	return yaml.NewEncoder(file).Encode(msg)
}

func (s *Storage) Update(msg dto.MessageTmpl) error {
	filePath := s.getFilePath(msg.Code)
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return storage.ErrMessageNotExists
	}
	file, fileErr := os.OpenFile(filePath, os.O_WRONLY, 0755)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	return yaml.NewEncoder(file).Encode(msg)
}

func (s *Storage) Rm(code string) error {
	return os.Remove(s.getFilePath(code))
}

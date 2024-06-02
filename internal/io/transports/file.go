package transports

import (
	"fmt"
	"os"
	"strings"

	"notifications/internal/dto"
)

type FileConfig struct {
	Path string
}

type File struct {
	cfg FileConfig
}

func (c *File) Send(to []string, msg *dto.Message, meta map[string]string) error {
	tmpl := `
------------------------------------------------------
	To: %s
	Meta: %s

	Subject: %s

	Body: %s
------------------------------------------------------
`
	mString := strings.Builder{}
	for mK, mV := range meta {
		mString.WriteString(fmt.Sprintf("%s:%s ", mK, mV))
	}

	f, fErr := os.OpenFile(c.cfg.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if fErr != nil {
		return fErr
	}
	defer f.Close()

	_, wErr := f.WriteString(fmt.Sprintf(tmpl, strings.Join(to, ","), mString.String(), msg.Subject, msg.Body))

	return wErr
}

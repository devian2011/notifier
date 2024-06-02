package transports

import (
	"errors"
	"fmt"

	"notifications/internal/dto"
)

type Transport interface {
	Send(to []string, msg *dto.Message, meta map[string]string) error
}

type Collection struct {
	transports map[string]Transport
}

func NewCollection(filePath string) (*Collection, error) {
	cfg, loadCfgErr := loadConfigFromFile(filePath)

	if loadCfgErr != nil {
		return nil, loadCfgErr
	}

	transports := map[string]Transport{}

	for fName, fCfg := range cfg.File {
		transports[fName] = &File{fCfg}
	}

	for smtpName, smtpCfg := range cfg.Smtp {
		transports[smtpName] = &SmtpSender{cfg: &smtpCfg}
	}

	return &Collection{transports: transports}, nil
}

func (c *Collection) AddTransport(code string, transport Transport) {
	c.transports[code] = transport
}

func (c *Collection) List() []string {
	result := make([]string, 0, len(c.transports))
	for name := range c.transports {
		result = append(result, name)
	}

	return result
}

func (c *Collection) Send(code string, to []string, msg *dto.Message, meta map[string]string) error {
	if transport, exists := c.transports[code]; exists {
		return transport.Send(to, msg, meta)
	} else {
		return errors.New(fmt.Sprintf("unknown transport with code %s", code))
	}
}

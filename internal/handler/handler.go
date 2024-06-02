package handler

import (
	"encoding/json"
	"errors"
	"strings"

	"notifications/internal/dto"
)

type Command string

const (
	Send Command = "send"

	List Command = "list"
	Get  Command = "get"

	Create Command = "create"
	Update Command = "update"
	Rm     Command = "remove"

	Transports Command = "transports"
)

type TransportCollection interface {
	List() []string
	Send(code string, to []string, msg *dto.Message, meta map[string]string) error
}

type Handler struct {
	storage    TemplateStorage
	transports TransportCollection
	sender     sender
}

func NewHandler(storage TemplateStorage, transports TransportCollection) *Handler {
	return &Handler{
		storage:    storage,
		transports: transports,
		sender: sender{
			storage:    storage,
			transports: transports,
		},
	}
}

func (h *Handler) Handle(command []byte, body []byte) (interface{}, error) {
	commandParts := strings.Split(string(command), "/")
	if len(commandParts) <= 0 {
		return nil, errors.New("unknown command")
	}
	action := Command(commandParts[len(commandParts)-1])
	switch action {
	case Send:
		rq := dto.MessageSendRequest{}
		uErr := json.Unmarshal(body, &rq)
		if uErr != nil {
			return nil, uErr
		}

		return h.sender.execute(rq)
	case Create:
		rq := dto.MessageTmpl{}
		uErr := json.Unmarshal(body, &rq)
		if uErr != nil {
			return nil, uErr
		}
		return rq, h.storage.Create(rq)
	case Update:
		rq := dto.MessageTmpl{}
		uErr := json.Unmarshal(body, &rq)
		if uErr != nil {
			return nil, uErr
		}
		return rq, h.storage.Update(rq)
	case Rm:
		rq := dto.MessageRequest{}
		uErr := json.Unmarshal(body, &rq)
		if uErr != nil {
			return nil, uErr
		}
		rmErr := h.storage.Rm(rq.Code)
		return nil, rmErr
	case List:
		return h.storage.List()
	case Get:
		rq := dto.MessageRequest{}
		uErr := json.Unmarshal(body, &rq)
		if uErr != nil {
			return nil, uErr
		}
		return h.storage.Load(rq.Code)
	case Transports:
		return h.transports.List(), nil
	default:
		return struct{ Command Command }{Command: action}, errors.New("unknown command")
	}
}

package handler

import (
	"testing"

	"notifications/internal/dto"
)

func Test_BuildMessage(t *testing.T) {
	tmpl := dto.MessageTmpl{
		Code:        "test",
		Description: "some desc",
		Params: map[string]dto.Param{
			"first": {
				Options: []string{},
				Default: "Bye",
			},
			"second": {
				Options: []string{},
				Default: "World",
			},
		},
		SubjectTmpl: "Subject {{.first}}",
		BodyTmpl:    "Body {{.first}} {{.second}}",
	}

	msg, _ := buildMessage(tmpl, map[string]interface{}{
		"first": "hello",
	})

	if msg.Subject != "Subject hello" {
		t.Error("subject is wrong", msg.Subject)
	}

	if msg.Body != "Body hello World" {
		t.Error("body is wrong", msg.Subject)
	}
}

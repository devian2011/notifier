package handler

import (
	"bytes"

	"html/template"

	"notifications/internal/dto"
)

type TemplateStorage interface {
	Load(code string) (dto.MessageTmpl, error)

	Create(msg dto.MessageTmpl) error
	Update(msg dto.MessageTmpl) error
	Rm(code string) error

	List() ([]dto.MessageTmpl, error)
}

func buildMessage(tmpl dto.MessageTmpl, inParams map[string]interface{}) (dto.Message, error) {
	msg := dto.Message{}
	params := buildParams(tmpl.Params, inParams)

	subjectBts := bytes.Buffer{}
	bodyBts := bytes.Buffer{}

	sTmpl, sTmplErr := template.New(tmpl.Code + "_subject").Parse(tmpl.SubjectTmpl)
	if sTmplErr != nil {
		return msg, sTmplErr
	}
	sExecErr := sTmpl.Execute(&subjectBts, params)
	if sExecErr != nil {
		return msg, sExecErr
	}

	bTmpl, bTmplErr := template.New(tmpl.Code + "_body").Parse(tmpl.BodyTmpl)
	if bTmplErr != nil {
		return msg, bTmplErr
	}
	bExecErr := bTmpl.Execute(&bodyBts, params)
	if bExecErr != nil {
		return msg, bExecErr
	}

	msg.Subject = subjectBts.String()
	msg.Body = bodyBts.String()

	return msg, nil
}

func buildParams(tmplParams map[string]dto.Param, inParams map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(tmplParams))
	for pName, pValue := range tmplParams {
		if inPVal, exists := inParams[pName]; exists {
			result[pName] = inPVal
		} else {
			result[pName] = pValue.Default
		}
	}

	return result
}

package dto

type MessageRequest struct {
	Code string `json:"code"`
}

type MessageSendRequest struct {
	Messages []struct {
		To        []string               `json:"to"`
		Code      string                 `json:"code,omitempty"`
		Meta      map[string]string      `json:"meta,omitempty"`
		Params    map[string]interface{} `json:"params,omitempty"`
		Message   *Message               `json:"message,omitempty"`
		Transport string                 `json:"transport"`
	} `json:"messages"`
}

type Param struct {
	Options []string    `json:"options"`
	Default interface{} `json:"default"`
}

type MessageTmpl struct {
	Code        string           `json:"code" yaml:"code"`
	Description string           `json:"description" yaml:"description"`
	Params      map[string]Param `json:"params" yaml:"params"`
	SubjectTmpl string           `json:"subject" yaml:"subject"`
	BodyTmpl    string           `json:"body" yaml:"body"`
}

type Message struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

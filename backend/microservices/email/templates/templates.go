package templates

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Template struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

var templates map[string]Template

func LoadTemplates() error {
	file, err := os.Open("./templates/templates.json")
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &templates); err != nil {
		return err
	}

	return nil
}

func ApplyTemplate(templateName string, args map[string]string) (Template, error) {
	if err := LoadTemplates(); err != nil {
		return Template{}, err
	}

	res := Template{}

	template, exist := templates[templateName]
    if !exist {
        return res, errors.New("Template with that name do not exist")
    }

    msg := template.Message

    for key, value := range args {
        placeholder := fmt.Sprintf("{{.%s}}", key)
        msg = strings.ReplaceAll(msg, placeholder, value)
    }

    res.Message = msg
    res.Sender = template.Sender
    res.Subject = template.Subject

	return res, nil
}

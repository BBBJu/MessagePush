package utils

import (
	"bytes"
	"fmt"
	"messagePush/models"
	"text/template"
)

func GetContentAfterTemplate(templateData map[string]interface{}, myTemplate models.Template) string {
	tmpl, err := template.New(myTemplate.Name).Parse(myTemplate.Content)
	if err != nil {
		fmt.Printf("模板解析失败: %v", err)
	}
	// 渲染模板
	var result bytes.Buffer
	err = tmpl.Execute(&result, templateData)
	if err != nil {
		fmt.Printf("模板渲染失败: %v", err)
	}
	fmt.Println("渲染结果:", result.String())
	return result.String()
}

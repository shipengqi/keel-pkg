package tmplutil

import (
	"bytes"
	"os"
	"text/template"
)

func ReplaceString(templ, templName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New(templName).Parse(templ)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ReplaceFile(templfile string, data interface{}) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.ParseFiles(templfile)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return  buf.String(), nil
}

func ReplaceAndWriteFile(templfile, output string, data interface{}) error {
	bugStr, err := ReplaceFile(templfile, data)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(bugStr))
	if err != nil {
		return err
	}
	return  nil
}

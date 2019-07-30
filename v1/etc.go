package notifier

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

func executeTemplate(str string, data interface{}) (string, error) {
	tmpl, err := template.New("").Parse(str)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")

	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func mustExecuteTemplate(str string, data interface{}) string {
	ret, err := executeTemplate(str, data)
	if err != nil {
		panic(err)
	}
	return ret
}

func getProp(obj map[string]interface{}, path string) string {
	names := strings.Split(path, ".")
	for i, name := range names {
		if val, ok := obj[name]; ok {
			if i == len(names)-1 {
				return val.(string)
			}

			obj = val.(map[string]interface{})
		} else {
			return ""
		}
	}
	return ""
}

func sendNotification(url, msg string, platformType PlatformType) error {
	platform := platforms[platformType]
	var reader *strings.Reader

	if platform.HTTPMethod == http.MethodPost {
		reader = strings.NewReader(msg)
	} else if platform.HTTPMethod == http.MethodGet {

	} else {
		return fmt.Errorf("invalid http method: " + platform.HTTPMethod)
	}

	req, err := http.NewRequest(platform.HTTPMethod, url, reader)
	if err != nil {
		return fmt.Errorf("create request failed: " + err.Error())
	}

	for k, v := range platform.HTTPHeaders {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request failed: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request status: %v, %v", resp.StatusCode, resp.Status)
	}

	return nil
}

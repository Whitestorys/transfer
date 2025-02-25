package image

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var (
	CCBackend = new(CC)
)

type CC struct {
	picBed
}

type CCResp struct {
	Code   int64         `json:"code"`
	Errors int64         `json:"total_error"`
	Image  []CCImageItem `json:"success_image"`
}

type CCImageItem struct {
	URL string `json:"url"`
	Del string `json:"delete"`
}

func (s CC) linkExtractor(link string) string {
	matcher := regexp.MustCompile("i[0-9]/[0-9]{4}/[0-9]{2}/[0-9]{2}/[a-zA-Z0-9]+")
	return matcher.FindString(link)
}

func (s CC) linkBuilder(link string) string {
	getter := regexp.MustCompile("i[0-9]/[0-9]{4}/[0-9]{2}/[0-9]{2}/[a-zA-Z0-9]+")
	return "https://upload.cc/" + getter.FindString(link) + ".png"
}

func (s CC) Upload(data []byte) (string, error) {

	body, err := s.upload(data, "https://upload.cc/image_upload", "uploaded_file[]")
	if err != nil {
		return "", err
	}

	var r CCResp

	if err := json.Unmarshal(body, &r); err != nil {
		return "", err
	}

	if r.Errors > 0 {
		return "", fmt.Errorf(string(body))
	}

	return s.linkBuilder(r.Image[0].URL), nil
}

package notion

import (
	"fmt"
	"net/http"
	"regexp"
	"transfer/apis"
)

var (
	matcher = regexp.MustCompile("https://(www\\.notion\\.so/signed/https%3A%2F%2F)?s3-us-west-2\\.amazonaws\\.com.*")
)

func (b notion) DoDownload(link string, config apis.DownConfig) error {
	err := apis.DownloadFile(&apis.DownloaderConfig{
		Link:     link,
		Config:   config,
		Modifier: apis.AddHeaders,
	})
	if err != nil {
		return fmt.Errorf("download failed on %s, returns %s\n", link, err)
	}
	return nil
}

func (b notion) AddHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; U; Linux x86_64; zh-CN; rv:1.9.2.10) "+
		"Gecko/20100922 Ubuntu/10.10 (maverick) Firefox/3.6.10")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", acceptLang)
	req.Header.Set("cookie", fmt.Sprintf("token_v2=%v", b.token))
}

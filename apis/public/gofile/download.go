package gofile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"transfer/apis"
	"transfer/utils"
)

var (
	matcher = regexp.MustCompile("(https://)?gofile\\.io/(\\?c=|download/|d/)[0-9a-zA-Z]{6}")
	reg     = regexp.MustCompile("=[0-9a-zA-Z]{6}")
)

func (b goFile) LinkMatcher(v string) bool {
	return matcher.MatchString(v)
}

func (b goFile) DoDownload(link string, config apis.DownConfig) error {
	err := b.download(link, config)
	if err != nil {
		return fmt.Errorf("download failed on %s, returns %s\n", link, err)
	}
	return nil
}

func (b goFile) download(v string, config apis.DownConfig) error {
	fileID := reg.FindString(v)
	fileID = fileID[1:]
	fmt.Printf("selecting server..")
	end := utils.DotTicker()
	body, err := http.Get(getServer)
	if err != nil {
		return fmt.Errorf("request %s returns error: %v", getServer, err)
	}

	data, err := ioutil.ReadAll(body.Body)
	if err != nil {
		return fmt.Errorf("read body returns error: %v", err)
	}
	_ = body.Body.Close()

	var sevData respBody
	if err := json.Unmarshal(data, &sevData); err != nil {
		return fmt.Errorf("parse body returns error: %v", err)
	}
	*end <- struct{}{}
	fmt.Printf("%s\n", strings.TrimSpace(sevData.Data.Server))
	fmt.Printf("fetching ticket..")
	end = utils.DotTicker()
	body, err = http.Get(fmt.Sprintf("https://apiv2.gofile.io/getUpload?c=%s", fileID))
	if err != nil {
		log.Println(fmt.Errorf("request %s returns error: %v", getServer, err))
		return fmt.Errorf("request %s returns error: %v", getServer, err)
	}

	data, err = ioutil.ReadAll(body.Body)
	if err != nil {
		return fmt.Errorf("read body returns error: %v", err)
	}
	_ = body.Body.Close()

	if err := json.Unmarshal(data, &sevData); err != nil {
		return fmt.Errorf("parse body returns error: %v", err)
	}
	*end <- struct{}{}

	fmt.Printf("%+v", sevData)
	for _, item := range sevData.Data.Items {
		err := apis.DownloadFile(&apis.DownloaderConfig{
			Link:     item.Link,
			Config:   config,
			Modifier: apis.AddHeaders,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

var checkHttpClient http.Client

type Check struct {
	Threshold int
	Url       string
	Template  string
}

type ElasticSearchResult struct {
	Count int `json:"count"`
}

func LoadChecks(checkFilesPattern string) (checks map[string]Check, err error) {
	paths, err := filepath.Glob(checkFilesPattern)
	if err != nil {
		return nil, err
	}

	checks = make(map[string]Check)

	for _, path := range paths {
		check := Check{}

		log.Printf("Loading check file '%s'.\n", path)
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		log.Printf("Parsing check file '%s'.\n", path)
		err = yaml.Unmarshal([]byte(data), &check)
		if err != nil {
			return nil, err
		}

		url := fmt.Sprintf("checks/%s", strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
		log.Printf("Setting the check URL for '%s' to '%s'.\n", path, url)
		check.Url = url
		checks[check.Url] = check
	}

	return checks, nil
}

func InitHttpClient(timeout int) {
	checkHttpClient = http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
}

func (check *Check) RunCheck(context *gin.Context, config Config) {
	resp, err := checkHttpClient.Post(config.ElasticSearchCountURL, "application/json", strings.NewReader(check.Template))
	if err != nil {
		context.String(503, "Error in making the request: %v", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		context.String(503, "Error in reading the request body: %v", err)
		return
	}

	result := ElasticSearchResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		context.String(503, "Error in decoding the JSON response: %v", err)
		return
	}

	// log.Println("Check %s found %d instances of the search in question with alert threshold %d.", check.Url, result.Count, check.Threshold)
	if result.Count >= check.Threshold {
		context.String(500, "Check '%s' found %d results, over or equal to the threshold of %d.", check.Url, result.Count, check.Threshold)
	} else {
		context.String(200, "Check '%s' found %d results, under the threshold of %d.", check.Url, result.Count, check.Threshold)
	}
}

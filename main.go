package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

type NotifyRequest struct {
	Interests []string `json:"interests"`
	Icon      string   `json:"icon"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
}

func send() (string, error) {
	request := NotifyRequest{
		Interests: []string{"allen.ccccnm@gmail.com"},
		Icon:      "",
		Title:     "点外卖",
		Body:      "",
	}
	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	res, err := http.Post("https://gotification.herokuapp.com/notify", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	var result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		PubId   string `json:"pubId"`
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(content, &result); err != nil {
		return "", err
	}
	if !result.Success {
		return "", fmt.Errorf("notify failed: %s", result.Message)
	}
	return result.PubId, nil
}

func main() {
	c := cron.New(cron.WithSeconds())
	if _, err := c.AddFunc("CRON_TZ=Asia/Shanghai * * 11 * * *", func() {
		if pubId, err := send(); err != nil {
			log.Printf("[FAIL]%s\n", err.Error())
		} else {
			log.Printf("[SUCCESS]PublishId: %s\n", pubId)
		}
	}); err != nil {
		panic(err)
	}
	c.Run()
}

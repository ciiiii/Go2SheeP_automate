package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"time"

	"github.com/ciiiii/Go2SheeP_notification/pusher"
	"github.com/robfig/cron/v3"
)

func send() (string, error) {
	request := pusher.NotifyRequest{
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
	defer res.Body.Close()
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

func currentTime() string {
	return time.Now().Format("2006.01.02 15:04:05")
}

func main() {
	c := cron.New()
	if _, err := c.AddFunc("CRON_TZ=Asia/Shanghai 0 11 * * *", func() {
		if pubId, err := send(); err != nil {
			fmt.Printf("[FAIL]%s|%s\n", currentTime(), err.Error())
		} else {
			fmt.Printf("[SUCCESS]%s|PublishId: %s\n", currentTime(), pubId)
		}
	}); err != nil {
		panic(err)
	}
	c.Run()
}

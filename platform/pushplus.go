package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PushPlus struct {
}

func (PushPlus *PushPlus) Run(pushPlusToken string, title string, content string) {
	url := "http://www.pushplus.plus/send/"
	var dataMap = make(map[string]string)
	dataMap["token"] = pushPlusToken
	dataMap["title"] = title
	dataMap["content"] = content
	dataByte, _ := json.Marshal(dataMap)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(dataByte))
	req.Header.Add("Content-Type", "application/json")
	time.Sleep(1 * time.Second)
	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		fmt.Println(string(body))
	} else {
		fmt.Println("微信推送失败")
	}
}

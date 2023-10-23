package platform

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Bilibili struct {
}

func (Bilibili *Bilibili) signIn(cookie string) (string, error) {
	url := "https://api.live.bilibili.com/xlive/web-ucenter/v1/sign/DoSign"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Cookie", cookie)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var resMap map[string]interface{}
	json.Unmarshal(body, &resMap)
	code, ok := resMap["code"].(float64)
	if ok {
		if code == 0 {
			if text, ok := resMap["data"].(map[string]interface{})["text"].(string); ok {
				return text, nil
			}
		} else {
			if msg, ok := resMap["message"].(string); ok {
				return msg, nil
			}
		}
	}
	return "", errors.New("不知名错误")
}

func (Bilibili *Bilibili) Run(pushPlusToken string, cookie string) {
	var title = "B站直播签到"
	pushPlus := PushPlus{}
	res, err := Bilibili.signIn(cookie)
	if err != nil {
		pushPlus.Run(pushPlusToken, title, err.Error())
	} else {
		pushPlus.Run(pushPlusToken, title, res)
	}
}

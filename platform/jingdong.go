package platform

import (
	"encoding/json"
	"io"
	"net/http"
)

type JD struct {
}

func (JD *JD) dailySign(cookie string) (string, error) {
	url := "https://api.m.jd.com/client.action?functionId=signBeanAct&body=%7B%22fp%22%3A%22-1%22%2C%22shshshfp%22%3A%22-1%22%2C%22shshshfpa%22%3A%22-1%22%2C%22referUrl%22%3A%22-1%22%2C%22userAgent%22%3A%22-1%22%2C%22jda%22%3A%22-1%22%2C%22rnVersion%22%3A%223.9%22%7D&appid=ld&client=apple&clientVersion=10.0.4&networkType=wifi&osVersion=14.8.1"
	req, err := http.NewRequest("POST", url, nil)
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
	if err = json.Unmarshal(body, &resMap); err != nil {
		return "", err
	}
	var title, content string
	if resMap["code"].(string) != "0" {
		title = "签到失败"
		if content, exist := resMap["errorMessage"]; exist {
			return title + content.(string), nil
		} else if content, exist := resMap["message"]; exist {
			return title + content.(string), nil
		}
		return "不知名错误", nil
	}
	var dailyAward map[string]interface{}
	if value, ok := resMap["data"].(map[string]interface{})["dailyAward"]; ok {
		dailyAward = value.(map[string]interface{})
	} else if value, ok := resMap["data"].(map[string]interface{})["continuityAward"]; ok {
		dailyAward = value.(map[string]interface{})
	} else if value, ok := resMap["data"].(map[string]interface{})["newUserAward"]; ok {
		dailyAward = value.(map[string]interface{})
	}
	title = dailyAward["title"].(string)
	if value, ok :=  dailyAward["beanAward"]; ok {
		content = value.(map[string]interface{})["beanCount"].(string) + "个京豆"
	} else {
		todayAward := dailyAward["awardList"].([]interface{})[0].(map[string]interface{})
		content = todayAward["beanCount"].(string) + "个京豆"
	}
	return title + content, nil
}

func (JD *JD) Run(pushPlusToken string, cookie string) {
	content, err := JD.dailySign(cookie)
	PushPlus := PushPlus{}
	if err != nil {
		PushPlus.Run(pushPlusToken, "京东每日签到", err.Error())
	} else {
		PushPlus.Run(pushPlusToken, "京东每日签到", content)
	}
}

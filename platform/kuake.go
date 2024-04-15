package platform

import (
	"autoSign/config"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type KK struct{}

func sign(kkCookie string) (msg string, err error) {
	url := "https://drive-m.quark.cn/1/clouddrive/capacity/growth/info?pr=ucpro&fr=pc&uc_param_str="
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", kkCookie)
	req.Header.Set("User-Agent", config.UserAgent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var data map[string]interface{}
	json.Unmarshal(body, &data)
	isSign := data["data"].(map[string]interface{})["cap_sign"].(map[string]interface{})["sign_daily"].(bool)
	if isSign != false {
		a := data["data"].(map[string]interface{})["cap_sign"].(map[string]interface{})
		number := a["sign_daily_reward"].(float64) / (1024 * 1024)
		return "今天已经签到过了,签到奖励容量" + strconv.Itoa(int(number)) + "MB", nil
	}
	signUrl := "https://drive-m.quark.cn/1/clouddrive/capacity/growth/sign?pr=ucpro&fr=pc&uc_param_str="
	payload := strings.NewReader(`{` + "" + `"sign_cyclic":"True"` + "" + `}`)
	req, err = http.NewRequest("POST", signUrl, payload)
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", kkCookie)
	req.Header.Set("User-Agent", config.UserAgent)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body1, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var data1 map[string]interface{}
	json.Unmarshal(body1, &data1)
	reward := data1["data"].(map[string]interface{})["sign_daily_reward"].(float64) / 2048
	return "签到成功，今日签到奖励" + strconv.Itoa(int(reward)) + "MB", nil
}

func (KK *KK) Run(pushPlusToken, cookie string) {
	msg, err := sign(cookie)
	PushPlus := PushPlus{}
	if err != nil {
		PushPlus.Run(pushPlusToken, "夸克网盘每日签到", err.Error())
	} else {
		PushPlus.Run(pushPlusToken, "夸克网盘每日签到", msg)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func GetAccessToken(refreshToken string) (string, error) {
	url := "https://auth.aliyundrive.com/v2/account/token"
	var dataMap = make(map[string]string)
	dataMap["grant_type"] = "refresh_token"
	dataMap["refresh_token"] = refreshToken
	dataByte, _ := json.Marshal(dataMap)
	req, err := http.NewRequest("POST", url, bytes.NewReader(dataByte))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
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
	if accessToken, ok := resMap["access_token"].(string); ok {
		return accessToken, nil
	}
	return "", errors.New("refreshToken过期,请更改后重试")
}

func SignIn(accessToken string) (string, error) {
	url := "https://member.aliyundrive.com/v1/activity/sign_in_list"
	data := []byte(`{
		"_rx-s":"mobile"
	}`)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", accessToken)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var resMap map[string]interface{}
	json.Unmarshal(body, &resMap)
	res := resMap["result"].(map[string]interface{})["signInCount"].(float64)
	signInCount := strconv.FormatFloat(res, 'f', 0, 64)
	return signInCount, nil
}

func GetReward(accessToken string, signInCount string) (string, error) {
	url := "https://member.aliyundrive.com/v1/activity/sign_in_reward?_rx-s=mobile"
	var dataMap = make(map[string]string)
	dataMap["signInDay"] = signInCount
	dataByte, _ := json.Marshal(dataMap)
	req, err := http.NewRequest("POST", url, bytes.NewReader(dataByte))
	if err != nil {
		return "", nil
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}
	var resMap map[string]interface{}
	json.Unmarshal(body, &resMap)
	if reward, ok := resMap["result"].(map[string]interface{})["notice"].(string); ok {
		return reward, nil
	}
	return "", errors.New("获取奖励失败")
}

func WXPush(pushPlusToken string, content string) {
	url := "http://www.pushplus.plus/send/"
	var dataMap = make(map[string]string)
	dataMap["token"] = pushPlusToken
	dataMap["title"] = "阿里云盘自动签到"
	dataMap["content"] = content
	dataByte, _ := json.Marshal(dataMap)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(dataByte))
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		fmt.Println(string(body))
	} else {
		fmt.Println("微信推送失败")
	}
}

func QianDao(refreshToken string) (string, string, error) {
	accessToken, err := GetAccessToken(refreshToken)
	if err != nil {
		return "", "", nil
	}
	signInCount, err := SignIn(accessToken)
	if err != nil {
		return "", "", nil
	}
	reward, err := GetReward(accessToken, signInCount)
	if err != nil {
		return "", "", nil
	}
	return signInCount, reward, nil
}

func main() {
	args := os.Args
	refreshToken := args[1]
	pushPlusToken := args[2]
	var signInCount string
	var reward string
	var err error
	signInCount, reward, err = QianDao(refreshToken)
	if err != nil {
		if err.Error() == "refreshToken过期,请更改后重试" {
			WXPush(pushPlusToken, "refreshToken过期,请更改后重试")
			fmt.Println("refreshToken过期,请更改后重试")
		} else {
			for i := 0; i < 100; i++ {
				signInCount, reward, err = QianDao(refreshToken)
				if err == nil {
					content := "签到成功，你已经签到" + signInCount + "次,本次签到奖励————" + reward
					fmt.Println(content)
					WXPush(pushPlusToken, content)
					break
				}
			}
		}
	} else {
		content := "签到成功，你已经签到" + signInCount + "次,本次签到奖励————" + reward
		fmt.Println(content)
		WXPush(pushPlusToken, content)
	}
}

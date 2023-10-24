package platform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type AliCloudDisk struct {
}

func (AliCloudDisk *AliCloudDisk) getAccessToken(refreshToken string) (string, error) {
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

func (AliCloudDisk *AliCloudDisk) signIn(accessToken string) (string, error) {
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

func (AliCloudDisk *AliCloudDisk) getReward(accessToken string, signInCount string) (string, error) {
	url := "https://member.aliyundrive.com/v1/activity/sign_in_reward?_rx-s=mobile"
	var dataMap = make(map[string]string)
	dataMap["signInDay"] = signInCount
	dataByte, _ := json.Marshal(dataMap)
	req, err := http.NewRequest("POST", url, bytes.NewReader(dataByte))
	if err != nil {
		return "", err
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
		return "", err
	}
	var resMap map[string]interface{}
	json.Unmarshal(body, &resMap)
	if reward, ok := resMap["result"].(map[string]interface{})["notice"].(string); ok {
		return reward, nil
	}
	return "", errors.New("获取奖励失败")
}

func (AliCloudDisk *AliCloudDisk) qianDao(refreshToken string) (string, string, error) {
	accessToken, err := AliCloudDisk.getAccessToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	signInCount, err := AliCloudDisk.signIn(accessToken)
	if err != nil {
		return "", "", err
	}
	reward, err := AliCloudDisk.getReward(accessToken, signInCount)
	if err != nil {
		return "", "", err
	}
	return signInCount, reward, nil
}

func (AliCloudDisk *AliCloudDisk) Run(pushPlusToken string, refreshToken string) {
	var signInCount string
	var reward string
	var err error
	var pushplus = PushPlus{}
	var title = "阿里云盘自动签到"
	signInCount, reward, err = AliCloudDisk.qianDao(refreshToken)
	if err != nil {
		if err.Error() == "refreshToken过期,请更改后重试" {
			pushplus.Run(pushPlusToken, title, "refreshToken过期,请更改后重试")
			fmt.Println("refreshToken过期,请更改后重试")
		} else {
			for i := 0; i < 100; i++ {
				signInCount, reward, err = AliCloudDisk.qianDao(refreshToken)
				if err == nil {
					content := "签到成功，你已经签到" + signInCount + "次,本次签到奖励————" + reward
					fmt.Println(content)
					pushplus.Run(pushPlusToken, title, content)
					break
				}
			}
		}
	} else {
		content := "签到成功，你已经签到" + signInCount + "次,本次签到奖励————" + reward
		fmt.Println(content)
		pushplus.Run(pushPlusToken, title, content)
	}
}

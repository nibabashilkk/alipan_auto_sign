package main

import (
	"autoSign/config"
	"autoSign/platform"
	"strings"
)

func main() {
	pushPlusToken := config.ConfigInstance.PushPlusToken
	refreshTokens := config.ConfigInstance.RefreshToken
	bilibiliCookies := config.ConfigInstance.BilibiliCookie
	jdCookies := config.ConfigInstance.JdCookie
	if refreshTokens != "" {
		refreshTokenList := strings.Split(refreshTokens, ",")
		aliCloudDisk := &platform.AliCloudDisk{}
		for _, refreshToken := range refreshTokenList {
			aliCloudDisk.Run(pushPlusToken, refreshToken)
		}
	}
	if bilibiliCookies != "" {
		bilibiliCookieList := strings.Split(bilibiliCookies, ",")
		bilibili := &platform.Bilibili{}
		for _, bilibiliCookie := range bilibiliCookieList {
			bilibili.Run(pushPlusToken, bilibiliCookie)
		}
	}
	if jdCookies != "" {
		jdCookiesList := strings.Split(jdCookies, ",")
		jd := &platform.JD{}
		for _, value := range jdCookiesList {
			jd.Run(pushPlusToken, value)
		}
	}
}

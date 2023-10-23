package main

import (
	"autoSign/platform"
	"os"
)

func main() {
	args := os.Args
	pushPlusToken := args[1]
	refreshToken := args[2]
	bilibiliCookie := args[3]
	if refreshToken != "null" {
		aliCloudDisk := platform.AliCloudDisk{}
		aliCloudDisk.Run(pushPlusToken, refreshToken)
	}
	if bilibiliCookie != "null" {
		bilibili := platform.Bilibili{}
		bilibili.Run(pushPlusToken, bilibiliCookie)
	}

}

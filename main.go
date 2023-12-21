package main

import (
	"autoSign/platform"
	"os"
	"strings"
)

func main() {
	args := os.Args
	pushPlusToken := args[1]
	refreshTokens := args[2]
	if refreshTokens != "null" {
		refreshTokenList := strings.Split(refreshTokens, ",")
		aliCloudDisk := &platform.AliCloudDisk{}
		for _, refreshToken := range refreshTokenList {
			aliCloudDisk.Run(pushPlusToken, refreshToken)
		}
	}
}

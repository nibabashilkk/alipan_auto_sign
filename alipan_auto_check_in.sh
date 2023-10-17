#!/bin/bash
#需要每月更新一次
refresh_token=$1
##获取access_token
access_token_command="curl --location --request POST 'https://auth.aliyundrive.com/v2/account/token' --header 'Content-Type: application/json' --data '{
    \"grant_type\": \"refresh_token\",
    \"refresh_token\": \"$refresh_token\"
}' | jq '.access_token' | sed 's/\"//g'";
access_token=`eval $access_token_command`
if [ "$access_token" = "null" ];
  then echo "refresh_token失效,请更换"
else
sleep 2
##自动签到
header="Authorization: $access_token";
sign_command="curl --location --request POST 'https://member.aliyundrive.com/v1/activity/sign_in_list' --header '$header' --header 'Content-Type: application/json' --data '{\"_rx-s\":\"mobile\"}'"
sign_in_count=`eval $sign_command | jq '.result.signInCount'`
sleep 2
##领取奖励
reward_command="curl --location --request POST 'https://member.aliyundrive.com/v1/activity/sign_in_reward?_rx-s=mobile' \
--header '$header' \
--header 'Content-Type: application/json' \
--data '{
    \"signInDay\": \"$sign_in_count\"
}'"
notice=`eval $reward_command | jq '.result.notice'`
echo "签到成功，你已经签到$sign_in_count次，本次签到奖励————$notice"
fi

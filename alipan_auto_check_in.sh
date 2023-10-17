#!/bin/bash
#需要每月更新一次
refresh_token=$1
#pushplus_token
pushplus_token=$2
##获取access_token
echo $refresh_token
access_token_command="curl --location --request POST 'https://auth.aliyundrive.com/v2/account/token' --header 'Content-Type: application/json' --data '{
    \"grant_type\": \"refresh_token\",
    \"refresh_token\": \"$refresh_token\"
}' | jq '.access_token' | sed 's/\"//g'";
access_token=`eval $access_token_command`
if [ "$access_token" = "null" ];
  then 
    pushplus_command="curl --location --request POST 'http://www.pushplus.plus/send/' \
--header 'Content-Type: application/json' \
--data '{
    \"token\":\"$pushplus_token\",
    \"title\":\"阿里云盘自动签到\",
    \"content\":\"refresh_token失效,请及时更换\"
}'"
    eval $pushplus_command;
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
notice=`eval $reward_command | jq '.result.notice' | sed 's/\"//g'`
pushplus_command="curl --location --request POST 'http://www.pushplus.plus/send/' \
--header 'Content-Type: application/json' \
--data '{
    \"token\":\"$pushplus_token\",
    \"title\":\"阿里云盘自动签到\",
    \"content\":\"签到成功，你已经签到$sign_in_count次，本次签到奖励————$notice\"
}'"
echo
eval $pushplus_command
fi
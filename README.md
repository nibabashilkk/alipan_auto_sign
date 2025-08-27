## 阿里云盘自动签到

2023-11-02更新：不需要配置仓库secret，直接修改项目里的**config.yaml**配置文件即可

2023-11-01更新：增加京东每日签到领豆，具体操作请看https://www.xiaoliu.life/p/20231101a

2023-10-25更新：增加B站直播签到，具体操作请看https://www.xiaoliu.life/p/20231025a


相较于昨天的实现，今天下班后加了下`pushplus`微信推送，能够实时了解每天签到的情况方便refresh_token失效的时候我们更换。

### pushplus

pushplus是一个专门的推送平台，利用服务号能主动给用户发消息的机制来进行推送的。

官网：[pushplus(推送加)-微信消息推送平台](https://www.pushplus.plus/)

![1](https://cdn.xiaoliu.life/tc/20231017a/1.webp)

使用的话需要登录注册且关注它的公众号后才能接收到消息推送。

![2](https://cdn.xiaoliu.life/tc/20231017a/2.webp)

登陆后会看见你自己的token，有了token后才能调用接口通过服务号给你自己发消息。

### github action

项目地址：[nibabashilkk/alipan_auto_sign: 阿里云盘每日自动签到shell脚本 (github.com)](https://github.com/nibabashilkk/alipan_auto_sign)

这个是我仓库的地址，需要你fork到自己仓库里面。

![3](https://cdn.xiaoliu.life/tc/20231017a/3.webp)

fork完后到你刚刚fork的仓库创建`repository secret`，具体路径是`Settings`->`Secrets and variables`->`Actions`->`New repository secret`。

![4](https://cdn.xiaoliu.life/tc/20231017a/4.webp)

注意新建的两个secret名字一定要是`refresh_token`和`pushplus_token`。

![5](https://cdn.xiaoliu.life/tc/20231017a/5.webp)

两个密钥都填好后是下面这个样子。

![6](https://cdn.xiaoliu.life/tc/20231017a/6.webp)

`refresh_token`的值是登陆后从阿里云盘的localStorage获取的，不知道怎么获取可以看看这篇文章——[阿里云盘每天自动签到，可以领超级会员 (qq.com)](https://mp.weixin.qq.com/s?__biz=Mzk0ODQwNzk1NA==&mid=2247489039&idx=1&sn=55c1d37978dfcb6f4f67cdaad0dc3b35&chksm=c3694df2f41ec4e43dd6a6a658ff9192b6014fad79d7bac25d32251efe6a443f50c1d7adb4a0&token=635348881&lang=zh_CN#rd)。

`pushplus_token`则是第一步登录pushplus后显示的token。

两个参数都弄好后应该就可以运行了，可以测试下配置是否正确。

![7](https://cdn.xiaoliu.life/tc/20231017a/7.webp)

点击后会立刻执行脚本，可以点进去运行任务里查看日志。

![8](https://cdn.xiaoliu.life/tc/20231017a/8.webp)

由于我微信推送太多次了，这次就没退过去，正常情况下是像下面那样能在公众号收到消息的。

![9](https://cdn.xiaoliu.life/tc/20231017a/9.webp)

如果`refresh_token`失效的话也会推送给你提醒你更换token。

![10](https://cdn.xiaoliu.life/tc/20231017a/10.webp)

refresh_token失效的话更新下secret即可。

![11](https://cdn.xiaoliu.life/tc/20231017a/11.webp)

默认是每天早上八点半运行，实际可能会到中午十二点才会运行（之前试过天气推送就会延迟）。

### 特别感谢

"本项目 CDN 加速及安全防护由 Tencent EdgeOne 赞助：EdgeOne 提供长期有效的免费套餐，包含不限量的流量和请求，覆盖中国大陆节点，且无任何超额收费，感兴趣的朋友可以点击下面的链接领取"

[亚洲最佳CDN、边缘和安全解决方案 - Tencent EdgeOne](https://edgeone.ai/?from=github)
![](https://edgeone.ai/media/34fe3a45-492d-4ea4-ae5d-ea1087ca7b4b.png)

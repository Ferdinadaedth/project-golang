# 问答社区后端项目 

## 1.语言和框架 :

本项目使用golang进行编写,以gin作为web框架。
## 2.功能实现:

* 注册登录:注册登录的数据以bcrypt进行加密后储存在mysql中,保证密码安全性，并且在登录成功生成jwt用于验证登录状态以及实现权限控制。并使用cors中间件解决跨域问题。

* 提问回答:已登录用户可以正常公开发问和回答，在每个问题的评论下实现楼中楼，在个人主页可编辑或删除自己的问题和回答，删除问题时会删除对应所有回答。

* 私信:用户也可向发问用户发起私信，用户可在用户主页面进行处理自己的私信或者回复别人的私信。

* 点赞:利用redis数据库的响应速度的优势以及非关系性实现点赞功能。

* 缓存:利用redis数据库实现对所有问题拉取后利用redis设定一个缓存，并在下一次拉取时优先查看缓存，提高用户体验。

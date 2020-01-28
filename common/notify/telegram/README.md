### 如何使用telegram bot来接收通知
- 创建telegram bot
   点击[telegram bot father](https://t.me/botfather),发送指令`/newbot`,然后会提示让输入bot的名称，发送你想创建bot的name，注意这个名字是bot显示的名称，发送之后还会提示让发送一个bot的username，这个username必须是不能重复的，因为使用这个名字才可以关注到这个bot，发送成功之后会给你发送一个`token`
- 如何发送给指定用户消息
	使用tg打开上一步创建的机器人，发送`/start`后会收到一个id，发送消息时将这个ID填写至Send的第一个参数


- Example
```go
func Send() {
    token := "929493383:AA..."
	telegram, err := NewTelegram(token)
	if err != nil {
		t.Fatal(err)
	}
	telegram.Send([]string{"chatid"}, "测试标题", "测试内容")

	time.Sleep(time.Second * 2)
}
```
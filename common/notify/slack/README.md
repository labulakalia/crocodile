### 如何使用slack channel来接收通知

- 添加一个来接收通知的channel 
  点击最左边的`Channel`旁边的`加号`,然后输入一个你想创建的channel的名称，点击`Create`，这时一个channel就创建好了  
  然后点击slack pc端左下角`Add more apps`，然后在搜索框输入`Imcoming WebHooks`,点击`Add`，这时会在浏览器打开一个新的页面，再次点击`Add To Slack`会进入`Imcoming WebHooks`的配置页面，然后下面会出现一个`Post to Channel`，并且还有一个选择框，然后点击`Choose a channel`，然后选择刚才创建的channel，点击下面的蓝色的按钮完成添加，保存`webhook URL`


### Example
```go
var (
	webhook = "https://hooks.slack.com/services/TGM152H5E/BSXFZALEB/sc..."
)
func Send() {
	slack := NewSlack(webhook)
	err := slack.Send([]string{"labulakalia"}, "测试标题", "测试内容")
	if err != nil {
		t.Error(err)
	}
}
```
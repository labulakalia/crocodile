### 如何使用钉钉群的机器人来接收通知
- 添加群机器人
   点击 **钉钉群设置**->**智能群助手**->**添加机器人**，选择**自定义**，然后点击**添加**，  
   填写**机器人名称**，**安全设置**，然后选择安全设置，有三类自定义关键词、加签、IP地址，任选一项然后点击完成  
   得到一个Webhook，
- 初始化
   使用`NewDing`来进行初始化，第一个参数就是这个Webhook url，第二个参数是安全设置,`CustomKey`是对应的的自定义关键字，`Sign`是加签，`IPCidr`是IP地址，注意如果选择了安全设置，这时会得到一个签名，然后把这个签名作为第三个参数传入，如果安全设置不是选择的加签，则第三个参数为空字符串即可。注意如果想在消息通知中@某个人则需要把这个人加入群众并且tos参数要穿入此人注册钉钉的手机号

### Example
```go
var (
secret := "SEC..."
webhook := "https://oapi.dingtalk.com/robot/send?access_token=..."
)
func SendDing() {
   ding := NewDing(webhook, Sign, secret)
   err := ding.Send([]string{"..."}, "测试标题", "测试内容")
   if err != nil {
      t.Error(err)
   }
}
```
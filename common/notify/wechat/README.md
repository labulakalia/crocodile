### 如何使用企业微信号来接收通知
- [注册企业微信](https://work.weixin.qq.com/wework_admin/register_wx?from=myhome)  
- 获取corpid  
  每个企业都拥有唯一的corpid，获取此信息可在管理后台**我的企业**－>**企业信息**->**企业ID**（需要有管理员权限）  
- 创建应用  
  点击**应用管理**->**自建**->**创建应用**，填写应用名称，添加成员，然后点击创建应用，然后点击进入新创建的应用，将**AgentId**和**Secret**两个参数记录下来 
- 在微信中接收通知  
为了在微信中可以直接接收到消息，需要微信扫码关注微工作台，点击 **我的企业**->**微工作台**->**邀请关注**，使用微信关注即可接收到通知

### Example
```go
var (
   corpid = "wwb..."
   agentID = 1000002
   secret = "NgYcbPHa6DhR..."
)
func Send() {
	client := NewWeChat(corpid, agentID, secret)
	err := client.Send([]string{"..."}, "测试标题", "测试消息的文本")
	if err != nil {
		t.Error(err)
	}
}
```

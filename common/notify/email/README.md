### 邮箱报警
如使用自建邮件系统请设置 skipVerify 为 true 以避免证书校验错误

### Example
- SendEmailByTls
```go
func SendTLS() {
	address := "smtp.163.com"
	from := "...@163.com"
	password := "..."
	// if use qq
	// https://service.mail.qq.com/cgi-bin/help?subtype=1&&no=1001256&&id=28
	address = "smtp.qq.com"
	from = "...@qq.com"
	password = "..."

	to := "to1@email.com"
	subject:= ""
	body:=""

	tls := true
	anonymous := false
	// 如使用自建邮件系统请设置 skipVerify 为 true 以避免证书校验错误
	skipVerify = true
	port := 465
	s := NewSMTP(address, from, password, from, html, tls, anonymous, skipVerify, port)
	err := s.Send([]string{to}, subject, body)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

```

- SendEmail
```go
func Send() {
	address := "smtp.163.com"
	from := "...@163.com"
	password := "..."

	to := "to1@email.com"
	subject:= ""
	body:=""
	
	tls := true
	anonymous := false
	// 如使用自建邮件系统请设置 skipVerify 为 true 以避免证书校验错误
	skipVerify = true
	port := 25

	s := NewSMTP(address, from, password, from, plain, tls, anonymous, skipVerify, port)
	s.Send([]string{tos}, subject, body)
}
```

- SendEmailAnonyMous
```go
func SendAnonyMous() {
	address := "smtp.custom.com"
	from := "noreply@custom.com"
		
	to := "to1@email.com"
	subject:= ""
	body:=""
	
	tls := true
	anonymous := true
	// 如使用自建邮件系统请设置 skipVerify 为 true 以避免证书校验错误
	skipVerify = true
	port := 25

	s := NewSMTP(address, "", "", from, plain, tls, anonymous, skipVerify, 25)
	s.Send([]string{tos}, subject, body)
}
```
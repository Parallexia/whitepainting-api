package passportv2

import (
	"gopkg.in/gomail.v2"
)

//向邮箱发送激活账号邮件
func SendAccessEmail(email string) error {
	//发送的地址
	sendfrom := "example@example.com"
	//发送的页面 html
	page := ""

	m := gomail.NewMessage()
	m.SetHeader("From", sendfrom)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "账号邮箱验证")
	m.SetBody("text/html", page)

	d := gomail.NewDialer("smtp.example.com", 587, "user", "123456")

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

//向邮箱发送找回密码验证码邮件
func SendChangePasswdEmail(email string) error {
	//发送的地址
	sendfrom := "example@example.com"
	//发送的页面 html
	page := ""

	m := gomail.NewMessage()
	m.SetHeader("From", sendfrom)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "找回密码")
	m.SetBody("text/html", page)

	d := gomail.NewDialer("smtp.example.com", 587, "user", "123456")

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

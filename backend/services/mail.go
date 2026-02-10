package services

import (
	"crypto/rand"
	"e-woms/conf"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"std-library-slim/email"
	"std-library-slim/redis"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// 发生邮箱验证码
func SendEmailCode(email string) bool {
	_, err := redis.RDB().Exists(fmt.Sprintf(conf.KeyEmailValidCodeLock, email))
	if err != nil {
		logs.Error("[SendEmailCode][Get]Exists Redis Key KeyEmailValidCodeLock Error:", err, email)
		return false
	}
	code, _ := GenerateRandomNumberCode(6)
	if err := redis.RDB().Set(fmt.Sprintf(conf.KeyEmailValidCode, email), code, time.Duration(conf.KeyEmailValidCodeExpireTime)*time.Second); err != nil {
		logs.Error("[SendEmailCode][Get]Set Redis Key KeyEmailValidCode Error:", err, email)
		return false
	}
	if err := redis.RDB().Set(fmt.Sprintf(conf.KeyEmailValidCodeLock, email), code, time.Duration(conf.KeyEmailValidCodeLockExpireTime)*time.Second); err != nil {
		logs.Error("[SendEmailCode][Get]Set Redis Key KeyEmailValidCodeLock Error:", err, email)
		return false
	}

	body := fmt.Sprintf("您的验证码是: %s, 请在5分钟内使用", code)
	mail := OutLookEmail{}
	err = mail.Send(email, "", body)
	if err != nil {
		logs.Error("[SendEmailCode][Get] OutLookEmail.Sendcode:", code, email, err.Error())
		return false
	}
	logs.Info("[SendEmailCode][Get] KeyEmailValidCodeExpireTime, KeyEmailValidCodeLockExpireTime, code:", conf.KeyPhoneValidCodeExpireTime, conf.KeyPhoneValidCodeLockExpireTime, code, email)

	return true
}

// 发送通用邮箱
func SendCommonEmail(email, title, body string) error {
	mail := OutLookEmail{}
	err := mail.Send(email, title, body)
	if err != nil {
		logs.Error("[SendEmail][Get] OutLookEmail.Sendcode:", email, err.Error())
		return err
	}
	return nil
}

// 发送html漂亮文本
func SendCommonHTMLEmail(email, code string) error {
	// 发送邮件
	subject := ""
	content := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
	body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
	.container { max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f9f9f9; }
	.content { background-color: #ffffff; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
	.code { font-size: 32px; font-weight: bold; color: #ff5722; letter-spacing: 4px; text-align: center; padding: 20px; background-color: #f5f5f5; border-radius: 4px; margin: 20px 0; }
	.footer { margin-top: 20px; padding-top: 20px; border-top: 1px solid #eee; color: #999; font-size: 12px; }
</style>
</head>
<body>
<div class="container">
	<div class="content">
		<h2 style="color: #333; margin-top: 0;">🔐 验证码通知</h2>
		<p>您好！</p>
		<p>请使用以下验证码完成验证：</p>
		<div class="code">%s</div>
		<p>⚠️ <strong>重要提示：</strong></p>
		<ul>
			<li>验证码有效期为 <strong>5分钟</strong></li>
			<li>请勿将验证码透露给他人</li>
			<li>如非本人操作，请忽略此邮件</li>
		</ul>
		<div class="footer">
			<p>此邮件由系统自动发送，请勿回复。</p>
			<p>© 2024 官方系统</p>
		</div>
	</div>
</div>
</body>
</html>
`, code)

	err := SendCommonEmail(email, subject, content)
	if err != nil {
		logs.Error("[VerificationCode][Send] SendEmail error: %v", err)
		return err
	}

	return nil
}

type OutLookEmail struct{}

// SendMail 发送邮件帮助类
func (mail *OutLookEmail) Send(recipientEmail, title, body string) error {
	subjectTitle, _ := web.AppConfig.String("OUTLOOK_TITLE")
	senderEmail, _ := web.AppConfig.String("OUTLOOK_EMAIL")
	senderPassword, _ := web.AppConfig.String("OUTLOOK_PASSWORD")
	//smtpServer := "smtp.office365.com"    // 正确的 SMTP 服务器地址
	smtpServer := "smtp.gmail.com" // 正确的 SMTP 服务器地址

	subject := subjectTitle
	if title != "" {
		subject = title
	}

	port := "587" // 正确的端口（支持 STARTTLS）

	// 初始化邮件客户端
	email.New(&email.Option{
		Address:    smtpServer + ":" + port,
		AuthMethod: email.MethodPlainAuth,
		Auth:       email.Auth{Identity: "", Username: senderEmail, Password: senderPassword, Host: smtpServer},
	})

	// 构建支持HTML的邮件消息
	msg := mail.buildHTMLMessage(senderEmail, recipientEmail, subject, body)

	err := email.Cli().Send("no-reply", []string{recipientEmail}, msg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// buildHTMLMessage 构建HTML格式的邮件消息
func (mail *OutLookEmail) buildHTMLMessage(from, to, subject, htmlBody string) []byte {
	// 构建MIME格式的邮件头和正文
	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"

	// 组装邮件头
	var msg string
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n"

	// Base64编码HTML内容
	encoded := base64.StdEncoding.EncodeToString([]byte(htmlBody))

	// 每76个字符添加一个换行（符合RFC 2045标准）
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		msg += encoded[i:end] + "\r\n"
	}

	return []byte(msg)
}

func GenerateRandomNumberCode(length int) (string, error) {
	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[index.Int64()]
	}
	return string(result), nil
}

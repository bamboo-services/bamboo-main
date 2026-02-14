/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package test

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jordan-wright/email"
)

type smtpTestConfig struct {
	host      string
	port      int
	username  string
	password  string
	fromEmail string
	fromName  string
	toEmail   string
}

func loadSMTPTestConfig(t *testing.T) *smtpTestConfig {
	t.Helper()

	if testing.Short() {
		t.Skip("skip smtp e2e in short mode")
	}
	if os.Getenv("ENABLE_SMTP_E2E_TEST") != "true" {
		t.Skip("set ENABLE_SMTP_E2E_TEST=true to enable smtp e2e")
	}

	host := os.Getenv("EMAIL_HOST")
	portRaw := os.Getenv("EMAIL_PORT")
	username := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASS")
	fromEmail := os.Getenv("EMAIL_FROM")
	toEmail := os.Getenv("EMAIL_ADMIN_EMAIL")

	if host == "" || portRaw == "" || username == "" || password == "" || fromEmail == "" || toEmail == "" {
		t.Skip("EMAIL_HOST/EMAIL_PORT/EMAIL_USER/EMAIL_PASS/EMAIL_FROM/EMAIL_ADMIN_EMAIL are required")
	}

	port, err := strconv.Atoi(portRaw)
	if err != nil || port <= 0 {
		t.Skip("EMAIL_PORT is invalid")
	}

	fromName := os.Getenv("EMAIL_FROM_NAME")
	if fromName == "" {
		fromName = "竹叶"
	}

	return &smtpTestConfig{
		host:      host,
		port:      port,
		username:  username,
		password:  password,
		fromEmail: fromEmail,
		fromName:  fromName,
		toEmail:   toEmail,
	}
}

func sendMailWithTLS(t *testing.T, cfg *smtpTestConfig, subject string, html string) {
	t.Helper()

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", cfg.fromName, cfg.fromEmail)
	e.To = []string{cfg.toEmail}
	e.Subject = subject
	e.HTML = []byte(html)

	auth := smtp.PlainAuth("", cfg.username, cfg.password, cfg.host)
	addr := fmt.Sprintf("%s:%d", cfg.host, cfg.port)

	tlsConfig := &tls.Config{
		ServerName: cfg.host,
	}

	t.Logf("正在发送邮件到: %s", cfg.toEmail)
	t.Logf("SMTP 服务器: %s", addr)

	if err := e.SendWithTLS(addr, auth, tlsConfig); err != nil {
		t.Fatalf("邮件发送失败: %v", err)
	}
}

// TestSendEmail 测试发送邮件
//
// 使用 jordan-wright/email 库直接发送邮件，验证 SMTP 配置是否正确
func TestSendEmail(t *testing.T) {
	cfg := loadSMTPTestConfig(t)

	sendMailWithTLS(t, cfg, "【测试】Bamboo-Main 邮件模块测试", `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>邮件测试</title>
</head>
<body style="font-family: 'Microsoft YaHei', Arial, sans-serif; padding: 20px; background-color: #f5f5f5;">
    <div style="max-width: 600px; margin: 0 auto; background: #ffffff; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
        <h2 style="color: #333; border-bottom: 2px solid #4CAF50; padding-bottom: 10px;">
            🎉 邮件模块测试成功！
        </h2>
        <p style="color: #555; line-height: 1.8;">
            这是一封来自 <strong>Bamboo-Main</strong> 邮件模块的测试邮件。
        </p>
        <p style="color: #555; line-height: 1.8;">
            如果您收到了这封邮件，说明邮件发送功能已经配置正确！
        </p>
        <table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
            <tr>
                <td style="padding: 12px; border: 1px solid #ddd; background: #f9f9f9; width: 30%;">发送时间</td>
                <td style="padding: 12px; border: 1px solid #ddd;">`+time.Now().Format("2006-01-02 15:04:05")+`</td>
            </tr>
            <tr>
                <td style="padding: 12px; border: 1px solid #ddd; background: #f9f9f9;">SMTP 服务器</td>
                <td style="padding: 12px; border: 1px solid #ddd;">`+cfg.host+`</td>
            </tr>
            <tr>
                <td style="padding: 12px; border: 1px solid #ddd; background: #f9f9f9;">发件人</td>
                <td style="padding: 12px; border: 1px solid #ddd;">`+cfg.fromEmail+`</td>
            </tr>
        </table>
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0 20px;">
        <p style="color: #999; font-size: 12px; text-align: center;">
            此邮件由 Bamboo-Main 系统自动发送，请勿回复。
        </p>
    </div>
</body>
</html>
`)

	t.Log("✅ 邮件发送成功！")
}

// TestSendEmailWithPool 测试使用连接池发送邮件
func TestSendEmailWithPool(t *testing.T) {
	cfg := loadSMTPTestConfig(t)

	// 注意：端口 465 需要使用 TLS，email.NewPool 默认使用 STARTTLS
	sendMailWithTLS(t, cfg, "【测试】连接池邮件测试", `
<html>
<body>
    <h2>🚀 连接池测试</h2>
    <p>这是通过连接池发送的测试邮件。</p>
    <p>发送时间: `+time.Now().Format("2006-01-02 15:04:05")+`</p>
</body>
</html>
`)

	t.Log("✅ 连接池邮件发送成功！")
}

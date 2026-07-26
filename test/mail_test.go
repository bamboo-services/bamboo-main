/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	xEmail "github.com/bamboo-services/bamboo-base-go/plugins/email"
)

// TestSendEmail 基于 bamboo-base-go xEmail 插件测试邮件发送
//
// 启用方式:
//
//	ENABLE_SMTP_E2E_TEST=true go test ./test/ -run TestSendEmail -v
//
// 需配置 EMAIL_HOST/EMAIL_PORT/EMAIL_USER/EMAIL_PASS/EMAIL_FROM/EMAIL_TLS，
// 收件人取 EMAIL_ADMIN_EMAIL。
func TestSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skip smtp e2e in short mode")
	}
	if os.Getenv("ENABLE_SMTP_E2E_TEST") != "true" {
		t.Skip("set ENABLE_SMTP_E2E_TEST=true to enable smtp e2e")
	}
	if os.Getenv("EMAIL_HOST") == "" || os.Getenv("EMAIL_USER") == "" ||
		os.Getenv("EMAIL_PASS") == "" || os.Getenv("EMAIL_FROM") == "" ||
		os.Getenv("EMAIL_ADMIN_EMAIL") == "" {
		t.Skip("EMAIL_HOST/EMAIL_USER/EMAIL_PASS/EMAIL_FROM/EMAIL_ADMIN_EMAIL are required")
	}

	// 经框架 InitClient 从环境变量装配邮件客户端
	clientAny, err := xEmail.InitClient(context.Background())
	if err != nil {
		t.Fatalf("初始化邮件客户端失败: %v", err)
	}
	client, ok := clientAny.(*xEmail.EmailClient)
	if !ok {
		t.Fatalf("InitClient 返回类型断言失败: %T", clientAny)
	}

	to := os.Getenv("EMAIL_ADMIN_EMAIL")
	html := fmt.Sprintf(`
<html>
<body style="font-family: 'Microsoft YaHei', Arial, sans-serif; padding: 20px; background-color: #f5f5f5;">
    <div style="max-width: 600px; margin: 0 auto; background: #ffffff; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
        <h2 style="color: #333; border-bottom: 2px solid #4f7cff; padding-bottom: 10px;">
            🎉 xEmail 插件测试成功！
        </h2>
        <p style="color: #555; line-height: 1.8;">
            这是一封基于 <strong>bamboo-base-go xEmail 插件</strong> 发送的测试邮件。
        </p>
        <p style="color: #555; line-height: 1.8;">
            发送时间: %s
        </p>
    </div>
</body>
</html>
`, time.Now().Format("2006-01-02 15:04:05"))

	t.Logf("正在发送邮件到: %s", to)
	if err := client.Send(context.Background(), &xEmail.Message{
		To:       []string{to},
		Subject:  "【测试】Bamboo-Main xEmail 插件测试",
		HTMLBody: html,
	}); err != nil {
		t.Fatalf("邮件发送失败: %v", err)
	}

	t.Log("✅ 邮件发送成功！")
}

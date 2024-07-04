// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package mail

import (
	"context"

	"xiaoMain/api/mail/v1"
)

type IMailV1 interface {
	MailSend(ctx context.Context, req *v1.MailSendReq) (res *v1.MailSendRes, err error)
}

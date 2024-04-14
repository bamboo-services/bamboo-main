package test

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncode(t *testing.T) {
	ctx := new(context.Context)
	baseData := base64.StdEncoding.EncodeToString([]byte("123456"))
	glog.Info(*ctx, fmt.Sprint(baseData))
	password, err := bcrypt.GenerateFromPassword([]byte(baseData), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	glog.Info(*ctx, string(password))
	err2 := bcrypt.CompareHashAndPassword(password, []byte(base64.StdEncoding.EncodeToString([]byte("123456"))))
	if err2 != nil {
		glog.Info(*ctx, "密码错误")
		return
	}
	glog.Info(*ctx, "密码正确")
}

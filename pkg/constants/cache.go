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

package constants

import (
	"fmt"
	"strings"

	xEnv "github.com/bamboo-services/bamboo-base-go/defined/env"
)

type RedisKey string

const (
	RedisSystemConfig  RedisKey = "system:config:%s"
	RedisSystemUser    RedisKey = "system:user:%d"
	RedisAuthToken     RedisKey = "auth:token:%s"
	RedisEmailVerify   RedisKey = "email:verify:%s"
	RedisPasswordReset RedisKey = "password:reset:%s"
	RedisLinkFriend    RedisKey = "link:friend:%d"
	RedisLinkGroup     RedisKey = "link:group:%d"
	RedisLinkColor     RedisKey = "link:color:%d"
	RedisSponsorRecord RedisKey = "sponsor:record:%d"
	RedisSponsorChan   RedisKey = "sponsor:channel:%d"
)

func (k RedisKey) Get(args ...any) RedisKey {
	prefix := xEnv.GetEnvString(xEnv.NoSqlPrefix, "bm:")
	if !strings.HasSuffix(prefix, ":") {
		prefix += ":"
	}
	return RedisKey(fmt.Sprintf(prefix+string(k), args...))
}

func (k RedisKey) String() string {
	return string(k)
}

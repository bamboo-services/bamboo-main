package constants

import (
	"fmt"
	"strings"

	xEnv "github.com/bamboo-services/bamboo-base-go/env"
)

type RedisKey string

const (
	RedisSystemConfig  RedisKey = "system:config:%s"
	RedisSystemUser    RedisKey = "system:user:%d"
	RedisAuthToken     RedisKey = "auth:token:%s"
	RedisEmailVerify   RedisKey = "email:verify:%s"
	RedisPasswordReset RedisKey = "password:reset:%s"
	RedisMailQueue     RedisKey = "mail:queue"
	RedisMailRetry     RedisKey = "mail:retry"
	RedisMailFailed    RedisKey = "mail:failed:%s"
	RedisMailStats     RedisKey = "mail:stats"
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

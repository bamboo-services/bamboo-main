package logic

import (
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCache "github.com/bamboo-services/bamboo-base-go/major/cache"
	"gorm.io/gorm"
)

type logic struct {
	db    *gorm.DB
	cache *xCache.Manager
	log   *xLog.LogNamedLogger
}

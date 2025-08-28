package startup

import (
	"bamboo-main/internal/model/base"

	xConsts "github.com/bamboo-services/bamboo-base-go/constants"
	"github.com/mitchellh/mapstructure"
)

// ConfigInit 初始化并加载客制化配置文件。
//
// 此方法通过解码服务注册中的原始配置内容，将其转换为 `BambooConfig` 实例，并赋值给结构体的 `Config` 字段。
// 如果解码失败，程序将触发 panic 并终止运行。
//
// 注意: 使用此方法之前需确保服务注册中的配置已正确加载，以避免解码失败。
func (r *Reg) ConfigInit() {
	r.Serv.Logger.Named(xConsts.LogINIT).Info("客制化初始化配置行为")

	var getConfig base.BambooConfig
	err := mapstructure.Decode(r.Serv.Config, &getConfig)
	if err != nil {
		panic("[Config] 配置文件加载失败: " + err.Error())
	}

	r.Config = &getConfig
}

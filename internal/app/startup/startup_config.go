package startup

import (
	"context"
	"os"
	"strconv"

	"github.com/bamboo-services/bamboo-main/internal/models/base"
	bSdkConst "github.com/phalanx/beacon-sso-sdk/constant"

	xEnv "github.com/bamboo-services/bamboo-base-go/env"
	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xUtil "github.com/bamboo-services/bamboo-base-go/utility"
)

func (r *reg) configInit(ctx context.Context) (any, error) {
	log := xLog.WithName(xLog.NamedINIT)
	log.Info(ctx, "加载环境变量配置")

	nodeID := xEnv.GetEnvInt64(xEnv.SnowflakeNodeID, 1)
	businessCache := getEnvBoolByKey(bSdkConst.EnvSsoBusinessCache.String(), false)

	cfg := &base.BambooConfig{
		Xlf: base.BMConfig{
			Debug: xEnv.GetEnvBool(xEnv.Debug, false),
			Server: base.ServerConfig{
				Port: xEnv.GetEnvInt(xEnv.Port, 23333),
			},
		},
		Snowflake: base.SnowflakeConfig{
			NodeID: &nodeID,
		},
		Database: base.DatabaseConfig{
			Host:     xEnv.GetEnvString(xEnv.DatabaseHost, "localhost"),
			Port:     xEnv.GetEnvInt(xEnv.DatabasePort, 5432),
			User:     xEnv.GetEnvString(xEnv.DatabaseUser, "bamboo_main"),
			Pass:     xEnv.GetEnvString(xEnv.DatabasePass, ""),
			Name:     xEnv.GetEnvString(xEnv.DatabaseName, "bamboo_main"),
			Prefix:   xEnv.GetEnvString(xEnv.DatabasePrefix, "bm_"),
			SSLMode:  getEnvStringByKey("DATABASE_SSLMODE", "disable"),
			TimeZone: xEnv.GetEnvString(xEnv.DatabaseTimezone, "Asia/Shanghai"),
		},
		NoSQL: base.NoSQLConfig{
			Host:     xEnv.GetEnvString(xEnv.NoSqlHost, "localhost"),
			Port:     xEnv.GetEnvInt(xEnv.NoSqlPort, 6379),
			Pass:     xEnv.GetEnvString(xEnv.NoSqlPass, ""),
			Database: xEnv.GetEnvInt(xEnv.NoSqlDatabase, 0),
			Prefix:   xEnv.GetEnvString(xEnv.NoSqlPrefix, "bm"),
		},
		Email: base.EmailConfig{
			SMTPHost:    xEnv.GetEnvString(xEnv.EmailHost, ""),
			SMTPPort:    xEnv.GetEnvInt(xEnv.EmailPort, 465),
			Username:    xEnv.GetEnvString(xEnv.EmailUser, ""),
			Password:    xEnv.GetEnvString(xEnv.EmailPass, ""),
			FromEmail:   xEnv.GetEnvString(xEnv.EmailFrom, ""),
			FromName:    getEnvStringByKey("EMAIL_FROM_NAME", "竹叶"),
			AdminEmail:  getEnvStringByKey("EMAIL_ADMIN_EMAIL", ""),
			WorkerCount: getEnvIntByKey("EMAIL_WORKER_COUNT", 4),
			MaxRetry:    getEnvIntByKey("EMAIL_MAX_RETRY", 3),
			Timeout:     getEnvIntByKey("EMAIL_TIMEOUT", 10),
			UseTLS:      getEnvBoolByKey("EMAIL_USE_TLS", true),
			UseStartTLS: getEnvBoolByKey("EMAIL_USE_STARTTLS", false),
		},
		SSO: base.SSOConfig{
			ClientID:                 getEnvStringByKey(bSdkConst.EnvSsoClientID.String(), ""),
			ClientSecret:             getEnvStringByKey(bSdkConst.EnvSsoClientSecret.String(), ""),
			WellKnownURI:             getEnvStringByKey(bSdkConst.EnvSsoWellKnownURI.String(), ""),
			RedirectURI:              getEnvStringByKey(bSdkConst.EnvSsoRedirectURI.String(), ""),
			EndpointAuthURI:          getEnvStringByKey(bSdkConst.EnvSsoEndpointAuthURI.String(), ""),
			EndpointTokenURI:         getEnvStringByKey(bSdkConst.EnvSsoEndpointTokenURI.String(), ""),
			EndpointUserinfoURI:      getEnvStringByKey(bSdkConst.EnvSsoEndpointUserinfoURI.String(), ""),
			EndpointIntrospectionURI: getEnvStringByKey(bSdkConst.EnvSsoEndpointIntrospectionURI.String(), ""),
			EndpointRevocationURI:    getEnvStringByKey(bSdkConst.EnvSsoEndpointRevocationURI.String(), ""),
			BusinessCache:            &businessCache,
		},
	}

	return cfg, nil
}

func getEnvStringByKey(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	return value
}

func getEnvIntByKey(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getEnvBoolByKey(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	boolValue, ok := xUtil.Parse().Bool(value)
	if !ok {
		return defaultValue
	}
	return boolValue
}

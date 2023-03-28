package config

import (
	"time"

	"github.com/labstack/echo/v4/middleware"
)

func GetEchoLogConfig(appConf *AppConfig) middleware.LoggerConfig {
	echoLogConf := middleware.DefaultLoggerConfig
	echoLogConf.CustomTimeFormat = time.RFC3339
	// echoLogConf.Format = fmt.Sprintln(`{"level":"info","source":"echo","id":"${id}","mt":"${method}","uri":"${uri}","st":${status},"e":"${error}","lc":"${latency_human}","ts":"${time_custom}"}`)
	return echoLogConf
}

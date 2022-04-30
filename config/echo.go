package config

import (
	"fmt"

	"github.com/labstack/echo/v4/middleware"
)

func GetEchoLogConfig(appConf *AppConfig) middleware.LoggerConfig {
	echoLogConf := middleware.DefaultLoggerConfig
	echoLogConf.CustomTimeFormat = `2006-01-02T15:04:05.000Z0700`
	echoLogConf.Format = fmt.Sprintln(`{"level":"info","source":"echo","id":"${id}","mt":"${method}","uri":"${uri}","st":${status},"e":"${error}","lc":"${latency_human}","ts":"${time_custom}"}`)
	return echoLogConf
}

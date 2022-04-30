package utils

import (
	"reflect"
	"runtime"
	"time"

	"github.com/dzungtran/echo-rest-api/pkg/logger"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func ExecuteWithRetry(exeFunc func() (interface{}, bool, error), retries int, sleepTime time.Duration) (interface{}, error) {
	var err error
	var resp interface{}
	for i := 0; i < retries; i++ {
		var shouldRetry bool
		resp, shouldRetry, err = exeFunc()
		if !shouldRetry {
			return resp, err
		} else {
			logger.Log().Warnf("Error when execute func %v: %v, retry: %v", GetFunctionName(exeFunc), err, i)
			if sleepTime > 0 {
				time.Sleep(sleepTime)
			}
		}
	}

	return resp, err
}

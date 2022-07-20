package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log Logger
)

type (
	CallLogOption struct {
		applyFunc func(*logConfigs)
	}

	logConfigs struct {
		Level zapcore.Level
		// Encoding sets the logger's encoding. Valid values are "json" and "console"
		Encoding string
	}
	Logger interface {
		Infof(template string, args ...interface{})
		Infow(msg string, keysAndValues ...interface{})
		Info(args ...interface{})

		Debugf(template string, args ...interface{})
		Debugw(msg string, keysAndValues ...interface{})
		Debug(args ...interface{})

		Warnf(template string, args ...interface{})
		Warnw(msg string, keysAndValues ...interface{})
		Warn(args ...interface{})

		Errorf(template string, args ...interface{})
		Errorw(msg string, keysAndValues ...interface{})
		Error(args ...interface{})

		Panicf(template string, args ...interface{})
		Panicw(msg string, keysAndValues ...interface{})
		Panic(args ...interface{})

		Fatalf(template string, args ...interface{})
		Fatalw(msg string, keysAndValues ...interface{})
		Fatal(args ...interface{})

		Sync() error
	}
)

func init() {
	var cfg = zap.NewDevelopmentConfig()
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Encoding = "json"
	cfg.OutputPaths = []string{"stdout"}
	logger, _ := cfg.Build()
	log = logger.Sugar()
}

// InitLog override default config
func InitLog(env string) {
	if env == "production" {
		var cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.Encoding = "json"
		cfg.OutputPaths = []string{"stdout"}
		cfg.Level.SetLevel(zap.WarnLevel)
		logger, _ := cfg.Build()
		log = logger.Sugar()
	}
}

func Log() Logger {
	return log
}

func Set(newLog Logger) {
	log = newLog
}

func WithConfigLevel(level string) CallLogOption {
	return CallLogOption{
		applyFunc: func(oio *logConfigs) {
			switch strings.ToLower(level) {
			case "debug":
				oio.Level = zapcore.DebugLevel
			case "info":
				oio.Level = zapcore.InfoLevel
			case "warn":
				oio.Level = zapcore.WarnLevel
			case "error":
				oio.Level = zapcore.ErrorLevel
			case "fatal":
				oio.Level = zapcore.FatalLevel
			}
		},
	}
}

func WithConfigEncoding(encoding string) CallLogOption {
	return CallLogOption{
		applyFunc: func(oio *logConfigs) {
			oio.Encoding = strings.ToLower(encoding)
		},
	}
}

func InitWithOptions(cnf ...CallLogOption) {
	var cfg = zap.NewDevelopmentConfig()
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{"stdout"}
	cfg.Encoding = "json"

	c := &logConfigs{}
	if len(cnf) > 0 {
		c = applyLogConfig(cnf)
	}

	if c.Encoding != "" {
		cfg.Encoding = c.Encoding
	}

	if c.Level >= -1 {
		cfg.Level.SetLevel(zap.DebugLevel)
	}

	logger, _ := cfg.Build()

	log = logger.Sugar()
}

func applyLogConfig(callOptions []CallLogOption) *logConfigs {
	if len(callOptions) == 0 {
		return &logConfigs{}
	}

	optCopy := &logConfigs{}
	for _, f := range callOptions {
		f.applyFunc(optCopy)
	}
	return optCopy
}

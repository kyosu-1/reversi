package logger

import (
	stackdriver "github.com/tommy351/zap-stackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewStackdriverZapLogger(serviceName, version string, level zapcore.Level) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig = stackdriver.EncoderConfig

	return config.Build(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return &stackdriver.Core{
				Core: core,
			}
		}),
		zap.Fields(
			stackdriver.LogServiceContext(&stackdriver.ServiceContext{
				Service: serviceName,
				Version: version,
			})))
}

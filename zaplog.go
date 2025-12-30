package ptti

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initZap() *zap.Logger {
	encodecfg := zap.NewProductionEncoderConfig()
	encodecfg.TimeKey = "timestampt"
	encodecfg.EncodeTime = zapcore.ISO8601TimeEncoder

	env := GetEnv("ENV")

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       env != "production",
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encodecfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]any{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}

func Log() *zap.SugaredLogger {
	logger := initZap()
	defer logger.Sync()

	sugarredLog := logger.Sugar()

	return sugarredLog
}

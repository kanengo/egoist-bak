package zaplog

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	zapLogInitializer interface {
		logInit(config *zaplogConfig) (zap.AtomicLevel, *zap.Logger, error)
	}

	macZapLogInit struct {
	}

	winZapLogInit struct {
	}

	unixZapLogInit struct {
	}
)

func (*macZapLogInit) logInit(config *zaplogConfig) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		zapconfig zap.Config
		level     zap.AtomicLevel
		logger    *zap.Logger
		err       error
	)

	//if config.isDevelop {
	//	zapconfig = zap.NewDevelopmentConfig()
	//} else {
	//	zapconfig = zap.NewProductionConfig()
	//}
	zapconfig = zap.NewDevelopmentConfig()

	zapconfig.DisableStacktrace = true
	zapconfig.EncoderConfig.TimeKey = "timestamp"
	zapconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err = zapconfig.Build()
	level = zapconfig.Level

	return level, logger, err
}

func (*winZapLogInit) logInit(config *zaplogConfig) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		zapconfig zap.Config
		level     zap.AtomicLevel
		logger    *zap.Logger
		err       error
	)

	//if config.isDevelop {
	//	zapconfig = zap.NewDevelopmentConfig()
	//} else {
	//	zapconfig = zap.NewProductionConfig()
	//}

	zapconfig = zap.NewDevelopmentConfig()

	zapconfig.DisableStacktrace = true
	zapconfig.EncoderConfig.TimeKey = "timestamp"
	zapconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err = zapconfig.Build()
	level = zapconfig.Level

	return level, logger, err
}

func epochFullTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func (*unixZapLogInit) logInit(config *zaplogConfig) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		level  zap.AtomicLevel
		logger *zap.Logger
		err    error
	)

	writers := []zapcore.WriteSyncer{os.Stderr}
	output := zapcore.NewMultiWriteSyncer(writers...)
	if len(config.logPath) != 0 {
		output = zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.logPath,
			MaxSize:    config.maxSize,
			MaxAge:     config.maxAge,
			Compress:   config.compress,
			MaxBackups: config.maxBackups,
			LocalTime:  true,
		})
	}
	encConfig := zap.NewProductionEncoderConfig()
	encConfig.TimeKey = "timestamp"
	encConfig.EncodeTime = epochFullTimeEncoder

	encoder := zapcore.NewJSONEncoder(encConfig)
	if config.isDevelop {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger = zap.New(zapcore.NewCore(encoder, output, level), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return level, logger, err
}

func zapLogInit(config *zaplogConfig) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		logIniter zapLogInitializer
		level     zap.AtomicLevel
		logger    *zap.Logger
		err       error
	)

	if runtime.GOOS == "darwin" {
		logIniter = &macZapLogInit{}
	} else if runtime.GOOS == "windows" {
		logIniter = &winZapLogInit{}
	} else {
		logIniter = &unixZapLogInit{}
	}

	if level, logger, err = logIniter.logInit(config); err != nil {
		fmt.Printf("zapLogInit failed, err:%v\n", err)
		return level, logger, err
	}
	//
	//if config.withPid {
	//	llog = llog.With(zap.Int("pid", os.Getpid()))
	//}
	//

	if config.instance != "" {
		logger = logger.With(zap.String("instance", config.instance))
	}

	if config.service != "" {
		logger = logger.With(zap.String("service", config.service))
	}

	return level, logger, err
}

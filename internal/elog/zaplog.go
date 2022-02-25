package zaplog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

var zapLogger *zap.Logger
var loggerLevel zap.AtomicLevel

func init() {
	_, zapLogger, _ = zapLogInit(defaultConfig)
	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func InitLogger(options ...ConfigOption) error {
	var err error
	//var level zap.AtomicLevel
	config := defaultConfig

	for _, option := range options {
		option.apply(config)
	}

	loggerLevel, zapLogger, err = zapLogInit(config)

	if err != nil {
		fmt.Printf("zap log init fail, err:%v", err)
		return err
	}

	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(1))

	logLevelHttpServer(config, loggerLevel)
	return nil
}

func Debug(msg string, fields ...zap.Field) {
	zapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	zapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	zapLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	zapLogger.Error(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	zapLogger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	zapLogger.Panic(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	zapLogger.DPanic(msg, fields...)
}

func Sync() error {
	return zapLogger.Sync()
}

func SetLogLevel(level string) error {
	level = strings.ToLower(level)
	switch strings.ToLower(level) {
	case "debug", "info", "warn", "error", "fatal":
	case "all":
		level = "debug"
	case "off", "none":
		level = "fatal"
	default:
		return errors.New("not support level")
	}
	client := http.Client{}

	type payload struct {
		Level string `json:"level"`
	}

	myPayload := payload{Level: level}
	buf, err := json.Marshal(myPayload)
	if err != nil {
		return err
	}
	Info("SetLogLevel", zap.String("path", setLevelPath), zap.String("level", level))
	req, err := http.NewRequest("PUT", setLevelPath, bytes.NewReader(buf))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		Error("SetLogLevel failed", zap.Error(err), zap.String("path", setLevelPath))
		return err
	}

	defer resp.Body.Close()

	return nil
}

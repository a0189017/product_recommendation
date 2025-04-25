package logger

import (
	"os"
	"product_recommendation/pkg/config"
	"product_recommendation/pkg/types"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var once sync.Once

func new() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	config := config.GetConfig()

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)

	encoderConfig.LineEnding = "\n\n"

	var core zapcore.Core

	writeSyncers := []zapcore.WriteSyncer{
		zapcore.AddSync(os.Stdout),
	}
	enableRotateLog := config.Server.EnableRotateLog
	if enableRotateLog {
		hook := lumberjack.Logger{
			Filename:   config.LogFile.FileName,   // 日誌文件路徑
			MaxSize:    config.LogFile.MaxSize,    // 每個日誌文件保存的最大尺寸 單位：M
			MaxBackups: config.LogFile.MaxBackups, // 日誌文件最多保存多少個備份
			MaxAge:     config.LogFile.MaxAge,     // 文件最多保存多少天
			Compress:   true,                      // 是否壓縮
		}

		writeSyncers = append(writeSyncers, zapcore.AddSync(&hook))

	}
	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writeSyncers...),
		atomicLevel,
	)

	caller := zap.AddCaller()
	development := zap.Development()
	logger = zap.New(core, caller, development)
}

func convertMapToFields(objects []types.H) []zap.Field {
	fields := []zap.Field{}
	for _, obj := range objects {
		for k := range obj {
			fields = append(fields, zap.Any(k, obj[k]))
		}
	}
	return fields
}

func getInstance() *zap.Logger {
	once.Do(func() {
		new()
	})
	return logger
}

func Debug(msg string, obj ...types.H) {
	log(zap.DebugLevel, msg, obj...)
}

func Info(msg string, obj ...types.H) {
	log(zap.InfoLevel, msg, obj...)
}

func Warn(msg string, obj ...types.H) {
	log(zap.WarnLevel, msg, obj...)
}

func Error(msg string, obj ...types.H) {
	log(zap.ErrorLevel, msg, obj...)
}

func log(level zapcore.Level, msg string, objects ...types.H) {
	defer getInstance().Sync()
	switch level {
	case zap.DebugLevel:
		getInstance().Debug(msg, convertMapToFields(objects)...)
	case zap.InfoLevel:
		getInstance().Info(msg, convertMapToFields(objects)...)
	case zap.WarnLevel:
		getInstance().Warn(msg, convertMapToFields(objects)...)
	case zap.ErrorLevel:
		getInstance().Error(msg, convertMapToFields(objects)...)
	}
}

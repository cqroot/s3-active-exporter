package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLogger *zap.Logger

func getZapEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	return encoder
}

func Init(debug bool, caller bool) {
	writeSyncer := zapcore.AddSync(os.Stdout)

	var core zapcore.Core
	if debug {
		core = zapcore.NewCore(getZapEncoder(), writeSyncer, zapcore.DebugLevel)
	} else {
		core = zapcore.NewCore(getZapEncoder(), writeSyncer, zapcore.InfoLevel)
	}

	if caller {
		zapLogger = zap.New(core, zap.AddCaller())
	} else {
		zapLogger = zap.New(core)
	}
}

func Sync() {
	zapLogger.Sync()
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

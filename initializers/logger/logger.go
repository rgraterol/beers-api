package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"

	"github.com/rgraterol/beers-api/initializers/config"
)

// LoggerConfiguration represents configuration for logs.
type LoggerConfiguration struct {
	// Level of logging, can be DebugLevel, InfoLevel, WarnLevel, ErrorLevel, DPanicLevel, PanicLevel
	Level string `yaml:"level"`
}

var LoggerConfig LoggerConfiguration

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLogWriter() zapcore.WriteSyncer {
	path, err :=  os.Getwd()
	if err != nil {
		panic(err)
	}
	timeString := time.Now().Format("02-01-2006")
	file, err := os.OpenFile(path + "/logs/" + timeString +".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return zapcore.AddSync(file)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.UTC().Format("02-01-2006T15:04:05Z0700"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLevel() zapcore.Level {
	err := config.LoadConfigSection("logger", &LoggerConfig)
	if err != nil {
		panic(err)
	}
	if level, err := levelMap[LoggerConfig.Level]; err {
		return level
	}
	return zapcore.DebugLevel
}

func LoggerInitializer() {
	writerSync := getLogWriter()
	encoder := getEncoder()

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writerSync, getLevel()),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLevel()),
	)
	logg := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logg)
}
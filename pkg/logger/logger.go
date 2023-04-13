package logger

import (
	"log"
	"net"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(serviceName, logstashUrl, enviroment string) (*zap.Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		FunctionKey:    "logger_name",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	tcpConn, err := net.DialTimeout("tcp", logstashUrl, 5*time.Second)
	if err != nil {
		log.Printf("Failed to connect to logstash: %s\n", err)
		return nil, err
	}
	logstashConfig := zapcore.NewJSONEncoder(encoderConfig)
	logshtashSyncer := zapcore.AddSync(tcpConn)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(
			logstashConfig,
			logshtashSyncer,
			zap.NewAtomicLevelAt(zap.DebugLevel),
		),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)
	logstashConfig.AddString("environment", enviroment)
	logstashConfig.AddString("appname", serviceName)
	logstashConfig.AddString("Thread", "main")
	logstashConfig.AddString("hostip", getOutboundIP().String())
	logstashConfig.AddString("containerId", getHostname())

	logger := zap.New(core, zap.AddCaller(), zap.ErrorOutput(zapcore.Lock(zapcore.AddSync(tcpConn))))

	return logger, nil
}

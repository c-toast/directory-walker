package main

import (
	"github.com/c-toast/directory-walker/filewalker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.Logger

var moveDir = ""

func init() {
	//replace the default json decoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	//set up the output of logger
	writeSyncer := zapcore.AddSync(os.Stderr)

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	//logger = zap.New(core, zap.AddCaller())
	logger = zap.New(core)
}

func main() {
	originPath := ""
	moveDir = ""
	walker := filewalker.New()
	err := walker.Init(originPath)
	if err != nil {
		logger.Error("init error", zap.Error(err))
		return
	}

	walker.RegisterHandler(&unzipHandler{})
	walker.RegisterHandler(&removeHandler{})
	walker.RegisterHandler(&moveHandler{})

	walker.Walk()
	logger.Info("finish walking")
}

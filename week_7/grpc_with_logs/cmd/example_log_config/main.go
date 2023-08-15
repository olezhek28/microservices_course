package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	prodLog()
	sugarLog()
	samplingLog()
}

func prodLog() {
	logger := zap.Must(zap.NewProduction())

	defer logger.Sync()

	logger.Info("Hello from Zap logger!")

	logger.Info("Prod logger",
		zap.String("username", "oleg"),
		zap.Int("userid", 423232),
		zap.String("provider", "avl"),
	)
}

func sugarLog() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	logger.Sugar().Infow("Sugar logger", "1234", "userID")
}

func samplingLog() {
	stdout := zapcore.AddSync(os.Stdout)

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	productionCfg.StacktraceKey = "stack"

	jsonEncoder := zapcore.NewJSONEncoder(productionCfg)

	jsonOutCore := zapcore.NewCore(jsonEncoder, stdout, level)

	samplingCore := zapcore.NewSamplerWithOptions(
		jsonOutCore,
		time.Second, // interval
		3,           // log first 3 entries
		0,           // thereafter log zero entires within the interval
	)

	log := zap.New(samplingCore)

	for i := 1; i <= 10; i++ {
		log.Info("an info message")
		log.Warn("a warning")
	}
}

package logger

import "go.uber.org/zap"

// Build builds a logger with the specified log level.
func Build(level string) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	err := cfg.Level.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}
	cfg.DisableStacktrace = true

	logProd, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	log := logProd.Sugar()
	return log, nil
}

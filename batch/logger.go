package batch

import "go.uber.org/zap"

// NewLogger returns a new zap logger
func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	return sugar, nil
}

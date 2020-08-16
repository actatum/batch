package batch

import "go.uber.org/zap"

// NewLogger returns a new zap logger
func NewLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

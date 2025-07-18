package logger

import (
	"go.uber.org/zap"
)

func Init(production bool) (*zap.Logger, error) {
	if production {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

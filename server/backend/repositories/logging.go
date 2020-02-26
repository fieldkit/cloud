package repositories

import (
	"context"

	"go.uber.org/zap"

	"github.com/fieldkit/cloud/server/logging"
)

func Logger(ctx context.Context) *zap.Logger {
	return logging.Logger(ctx).Named("repos")
}

func OnlyLogIf(logger *zap.SugaredLogger, verbose bool) *zap.SugaredLogger {
	if !verbose {
		return zap.NewNop().Sugar()
	}
	return logger
}

package types

import (
	"context"
	"go.uber.org/zap"
)

type DebugLogger struct {
	context.Context
	*zap.SugaredLogger
}

func (l *DebugLogger) Log(values ...interface{}) {
	for _, item := range values {
		l.Debugf("%+v", item)
	}
}

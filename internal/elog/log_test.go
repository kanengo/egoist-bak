package zaplog

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	Debug("debug logger", zap.Time("mtime", time.Now()))
	Info("Info logger", zap.Time("mtime", time.Now()))
	Error("Error logger", zap.Time("mtime", time.Now()))
}

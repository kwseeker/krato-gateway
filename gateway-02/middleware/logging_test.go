package middleware

import (
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestKratosLog(t *testing.T) {

	logger := log.NewStdLogger(os.Stdout)
	l := log.NewHelper(logger)

	l.Log(log.LevelInfo, "foo", "bar")
	l.Debug()
	l.Info()
	l.Fatal()
}

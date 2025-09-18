package generator

import (
	"log"

	"github.com/kukymbr/dbmodgen/internal/util"
)

func newGennaLogger() *log.Logger {
	return log.New(&gennaLoggerWriter{}, "[genna] ", log.LstdFlags)
}

type gennaLoggerWriter struct{}

func (g *gennaLoggerWriter) Write(msg []byte) (n int, err error) {
	util.PrintDebugf(string(msg))

	return len(msg), nil
}

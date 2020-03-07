package helpers

import (
	"log"
	"os"
)

func CreateLogger(prefix string) *log.Logger {
	return log.New(os.Stderr, prefix+" ", log.LstdFlags+log.Lmsgprefix)
}

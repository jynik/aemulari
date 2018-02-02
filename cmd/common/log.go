package common

import (
	"github.com/op/go-logging"
	"os"
)

func InitLogging() {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logFmt := logging.MustStringFormatter("%{color}[%{level:.8s}]%{color:reset} %{message}")
	logFormatter := logging.NewBackendFormatter(logBackend, logFmt)
	logging.SetBackend(logFormatter)
}

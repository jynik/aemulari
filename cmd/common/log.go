package common

import (
	"github.com/op/go-logging"
	"os"
)

const fmt string = "%{color}[%{level:.8s}]%{color:reset} %{message}"

func InitLogging() {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logFmt := logging.MustStringFormatter(fmt)
	logFormatter := logging.NewBackendFormatter(logBackend, logFmt)
	logging.SetBackend(logFormatter)
}

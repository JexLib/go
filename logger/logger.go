package logger

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	logging "github.com/op/go-logging"
)

var Log *logging.Logger = nil
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.000} %{level:.4s} ▶ %{shortfunc} %{id:03x}%{color:reset} %{message}`,
)
var formatFile = logging.MustStringFormatter(
	`%{time:2006-01-02 15:04:05.000} %{level:.4s} ▶ %{shortfunc} %{id:03x} %{message}`,
)

func CreateLogger(isStdout bool, app ...string) {
	appstr := "application"
	if len(app) > 0 {
		appstr = app[0]
	}
	Log = logging.MustGetLogger(appstr)
	fullExeFilename, _ := exec.LookPath(os.Args[0])
	fullPath := filepath.Dir(fullExeFilename)
	logPath := "logs"
	if logFileWriter, er := NewRotateWriter(filepath.Join(fullPath, logPath), "log", ".log"); er == nil {

		beScreen := logging.NewLogBackend(os.Stdout, "", 0)
		beScreenFormatter := logging.NewBackendFormatter(beScreen, format)
		beScreenLeveled := logging.AddModuleLevel(beScreenFormatter)
		beScreenLeveled.SetLevel(logging.DEBUG, "")

		beFile := logging.NewLogBackend(logFileWriter, "", 0)
		beFileFormatter := logging.NewBackendFormatter(beFile, formatFile)
		beFileLeveled := logging.AddModuleLevel(beFileFormatter)
		beFileLeveled.SetLevel(logging.NOTICE, "")

		if isStdout {
			logging.SetBackend(beScreenLeveled, beFileLeveled)
		} else {
			logging.SetBackend(beFileLeveled)
		}
	} else {
		fmt.Errorf("Log init failure: %s", er)
	}

}

package logService

import logrusInfra "go-rotate-backup-s3/commons/infra/logrus"

func Info(msg string) {
	logrusInfra.Log.Info(msg)
}

func Error(msg string) {
	logrusInfra.Log.Error(msg)
}

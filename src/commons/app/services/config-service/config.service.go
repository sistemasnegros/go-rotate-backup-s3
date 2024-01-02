package configService

import (
	logService "go-rotate-backup-s3/commons/app/services/log-service"
	configInfra "go-rotate-backup-s3/commons/infra/godotenv"
	"strconv"
)

func GetJwtSecret() string {
	return configInfra.Get("JWT_SECRET")
}

func GetJwtExt() int {
	value, err := strconv.Atoi(configInfra.Get("JWT_EXP"))
	if err != nil {
		logService.Error(err.Error())
		return 14440
	}
	return value
}

func GetDbConfig() string {
	return configInfra.Get("DB_CONFIG")
}

func GetSmtpHost() string {
	return configInfra.Get("SMTP_HOST")
}

func GetSmtpPort() int {
	value, err := strconv.Atoi(configInfra.Get("SMTP_PORT"))
	if err != nil {
		panic("SMTP_PORT must be a int")
	}
	return value
}

func GetSmtpUser() string {
	return configInfra.Get("SMTP_USER")
}

func GetSmtpPass() string {
	return configInfra.Get("SMTP_PASS")
}

func GetSmtpFrom() string {
	return configInfra.Get("SMTP_FROM")
}

func GetSmtpTo() string {
	return configInfra.Get("SMTP_TO")
}

func GetSmtpEnabled() string {
	return configInfra.Get("SMTP_ENABLED")
}

func GetSmtpTemplate() string {
	return configInfra.Get("SMTP_TEMPLATE")
}

func GetMongoDb() string {
	return configInfra.Get("MONGO_DB")
}

type S3Config struct {
	AWS_SECRET_ACCESS_KEY string
	AWS_ACCESS_KEY_ID     string
	AWS_BUCKET            string
	AWS_REGION            string
	AWS_ENDPOINT          string
	AWS_URL_PUBLIC        string
}

func GetS3() *S3Config {
	return &S3Config{
		AWS_ACCESS_KEY_ID:     configInfra.Get("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS_KEY: configInfra.Get("AWS_SECRET_ACCESS_KEY"),
		AWS_BUCKET:            configInfra.Get("AWS_BUCKET"),
		AWS_REGION:            configInfra.Get("AWS_REGION"),
		AWS_ENDPOINT:          configInfra.Get("AWS_ENDPOINT"),
		AWS_URL_PUBLIC:        configInfra.Get("AWS_URL_PUBLIC"),
	}
}

type BackupConfig struct {
	BACKUP_KEEP            int
	BACKUP_SRC             string
	BACKUP_DST             string
	BACKUP_PREFiX_NAME     string
	BACKUP_COMMAND         string
	BACKUP_COMMAND_TIMEOUT int
}

func GetBackup() *BackupConfig {
	BACKUP_KEEP, err := strconv.Atoi(configInfra.Get("BACKUP_KEEP"))
	if err != nil {
		logService.Error(err.Error())
		panic(1)
	}

	BACKUP_COMMAND_TIMEOUT, err := strconv.Atoi(configInfra.Get("BACKUP_COMMAND_TIMEOUT"))
	if err != nil {
		logService.Error(err.Error())
		panic(1)
	}

	return &BackupConfig{
		BACKUP_KEEP:            BACKUP_KEEP,
		BACKUP_SRC:             configInfra.Get("BACKUP_SRC"),
		BACKUP_DST:             configInfra.Get("BACKUP_DST"),
		BACKUP_PREFiX_NAME:     configInfra.Get("BACKUP_PREFiX_NAME"),
		BACKUP_COMMAND:         configInfra.Get("BACKUP_COMMAND"),
		BACKUP_COMMAND_TIMEOUT: BACKUP_COMMAND_TIMEOUT,
	}
}

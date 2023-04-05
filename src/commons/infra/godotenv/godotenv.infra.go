package godotenvInfra

import (
	"flag"
	logService "go-rotate-backup-s3/commons/app/services/log-service"
	configDomain "go-rotate-backup-s3/commons/domain/config"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

func Load() {
	configPtr := flag.String("config", "../.env", "set define file config")
	flag.Parse()

	err := godotenv.Load(*configPtr)
	if err != nil {
		logService.Error("error loading .env")
		panic(1)
	}

	// Validate field
	config := &configDomain.Config{}
	val := reflect.ValueOf(config).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i).Name
		if Get(field) == "" {
			logService.Error("undefine in env: " + field)
			os.Exit(1)
		}
		logService.Info("load env: " + field)
	}

}

func Get(key string) string {
	value := os.Getenv(key)
	return value
}

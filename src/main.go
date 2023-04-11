package main

import (
	"context"
	smtpService "go-rotate-backup-s3/commons/app/services/smtp-service"
	msgDomain "go-rotate-backup-s3/commons/domain/msg"
	godotenvInfra "go-rotate-backup-s3/commons/infra/godotenv"
	gomailInfra "go-rotate-backup-s3/commons/infra/gomail"
	logrusInfra "go-rotate-backup-s3/commons/infra/logrus"
	s3Infra "go-rotate-backup-s3/commons/infra/s3"
	filesService "go-rotate-backup-s3/files/app/services"
	filesRepo "go-rotate-backup-s3/files/infra/s3"
	mainService "go-rotate-backup-s3/main/app/services"

	"go.uber.org/fx"
)

func main() {

	godotenvInfra.Load()
	logrusInfra.Init()
	msgDomain.New()

	app := fx.New(
		fx.NopLogger,

		fx.Provide(

			s3Infra.New,
			gomailInfra.New,
			smtpService.New,

			filesRepo.New,
			filesService.New,
			mainService.New,
		),
		fx.Invoke(
			start,
		),
	)

	app.Run()
}

func start(lc fx.Lifecycle, main *mainService.MainService, shutdowner fx.Shutdowner) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			go main.Run()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

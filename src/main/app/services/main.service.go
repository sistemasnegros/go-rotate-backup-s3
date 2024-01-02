package mainService

import (
	"context"
	"errors"
	"fmt"
	configService "go-rotate-backup-s3/commons/app/services/config-service"
	logService "go-rotate-backup-s3/commons/app/services/log-service"
	smtpService "go-rotate-backup-s3/commons/app/services/smtp-service"
	filesDomain "go-rotate-backup-s3/commons/domain/files"
	smtpDomain "go-rotate-backup-s3/commons/domain/smtp"
	"os"
	"os/exec"
	"strings"
	"time"

	filesService "go-rotate-backup-s3/files/app/services"

	"go.uber.org/fx"
)

type MainService struct {
	fileService *filesService.FilesService
	smtpService smtpService.ISmtpService
	smtpMessage smtpDomain.SendArgs
	shutdowner  fx.Shutdowner
}

func New(
	fileService *filesService.FilesService,
	smtpService smtpService.ISmtpService,
	shutdowner fx.Shutdowner,
) *MainService {
	return &MainService{
		fileService: fileService,
		smtpService: smtpService,
		shutdowner:  shutdowner,
		smtpMessage: smtpDomain.SendArgs{
			Template: "notify.html",
			Data: smtpDomain.EmailTemplateDefault{
				FullName:      "Devops",
				ButtonMessage: "Go to AWS",
				URL:           "https://s3.console.aws.amazon.com/s3/buckets/draketech-backups?region=us-east-2&bucketType=general&tab=objects",
			},
		},
	}
}

func (main *MainService) SendEmail() {

	if configService.GetSmtpEnabled() != "true" {
		logService.Info("notify smtp disabled")
		return
	}

	smtpTo := configService.GetSmtpTo()
	smtpToList := strings.Split(smtpTo, ",")

	for _, to := range smtpToList {
		main.smtpMessage.To = to
		errSmtp := main.smtpService.Send(main.smtpMessage)
		if errSmtp != nil {
			continue
		}

	}
}

func (main *MainService) SendEmailSuccess(message string) {
	main.smtpMessage.Subject = fmt.Sprintf("Backup mongodb successful: %s", configService.GetBackup().BACKUP_PREFiX_NAME)
	main.smtpMessage.Data.Message = message
	main.SendEmail()
}

func (main *MainService) SendEmailError(err string) {
	main.smtpMessage.Subject = fmt.Sprintf("Backup mongodb error: %s", configService.GetBackup().BACKUP_PREFiX_NAME)
	main.smtpMessage.Data.Message = fmt.Sprintf("Backup mongodb : %s", err)

	main.SendEmail()
}

func (main *MainService) handleError(err error) {
	if err != nil {
		logService.Error(err.Error())
		main.SendEmailError(err.Error())
		// panic(err)
		//main.shutdowner.Shutdown(fx.ExitCode(2))
		os.Exit(1)

	}
}

func (main *MainService) Run() {
	configBackup := configService.GetBackup()

	ctx := context.Background()
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(configBackup.BACKUP_COMMAND_TIMEOUT)*time.Second)
	defer cancel()

	// command := strings.Split(configBackup.BACKUP_COMMAND, " ")
	// cmd := exec.Command(command[0], command[1:]...)

	cmd := exec.CommandContext(ctx, "bash", "-c", configBackup.BACKUP_COMMAND)
	out, err := cmd.CombinedOutput()

	if err != nil {
		errorFormated := fmt.Sprintf("executed command failed: %s \n %s", err.Error(), string(out))
		main.handleError(errors.New(errorFormated))
	}

	logService.Info("executed command successful")

	file, err := os.Open(configBackup.BACKUP_SRC)

	main.handleError(err)

	defer file.Close()

	currentTime := time.Now()
	prefixName := fmt.Sprintf("%s_%s", currentTime.Format("2006-01-02_15:04:05"), configBackup.BACKUP_PREFiX_NAME)

	fileDomain := &filesDomain.File{
		Id:   configBackup.BACKUP_DST + "/v0/" + prefixName,
		Data: file,
		Name: prefixName,
	}

	rangeGenerated := make([]int, configBackup.BACKUP_KEEP)

	keep := configBackup.BACKUP_KEEP

	for range rangeGenerated {
		dst := fmt.Sprintf("%s/v%d/", configBackup.BACKUP_DST, keep)
		src := fmt.Sprintf("%s/v%d/", configBackup.BACKUP_DST, keep-1)

		srcListFolder, err := main.fileService.List(src)
		main.handleError(err)

		dstListFolder, err := main.fileService.List(dst)
		main.handleError(err)

		for _, v := range dstListFolder {
			err := main.fileService.Delete(v.Id)
			main.handleError(err)
		}

		for _, v := range srcListFolder {
			_, err := main.fileService.Copy(v.Id, dst+v.Name)
			main.handleError(err)

			err = main.fileService.Delete(v.Id)
			main.handleError(err)
		}

		keep = keep - 1
	}

	fileUpdated, err := main.fileService.UploadLarge(*fileDomain)

	main.handleError(err)

	logService.Info(fmt.Sprintf("create backup successful: %s", fileUpdated.Id))
	os.RemoveAll(configBackup.BACKUP_SRC)

	message := fmt.Sprintf("Backup mongodb successful: %s", "s3://"+configService.GetS3().AWS_BUCKET+"/"+configService.GetBackup().BACKUP_DST+"/v0/"+prefixName)
	main.SendEmailSuccess(message)

	main.shutdowner.Shutdown()
}

package mainService

import (
	"context"
	"fmt"
	configService "go-rotate-backup-s3/commons/app/services/config-service"
	logService "go-rotate-backup-s3/commons/app/services/log-service"
	filesDomain "go-rotate-backup-s3/commons/domain/files"
	"os"
	"os/exec"
	"time"

	filesService "go-rotate-backup-s3/files/app/services"

	"go.uber.org/fx"
)

type MainService struct {
	fileService *filesService.FilesService
	shutdowner  fx.Shutdowner
}

func New(fileService *filesService.FilesService, shutdowner fx.Shutdowner) *MainService {
	return &MainService{fileService: fileService, shutdowner: shutdowner}
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
		logService.Error("executed command failed")
		logService.Error(string(out))
		logService.Error(err.Error())
		panic(err)
	}

	logService.Info("executed command successful")

	file, err := os.Open(configBackup.BACKUP_SRC)

	if err != nil {
		fmt.Println("cannot able to read the file", err)
		logService.Error(err.Error())
		panic(1)
	}

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

		srcListFolder, _ := main.fileService.List(src)
		dstListFolder, _ := main.fileService.List(dst)

		for _, v := range dstListFolder {
			main.fileService.Delete(v.Id)
		}

		for _, v := range srcListFolder {
			main.fileService.Copy(v.Id, dst+v.Name)
			main.fileService.Delete(v.Id)
		}

		keep = keep - 1
	}

	fileUpdated, err := main.fileService.UploadLarge(*fileDomain)

	if err != nil {
		logService.Error(err.Error())
		panic("error in main")
	}

	logService.Info(fmt.Sprintf("create backup successful: %s", fileUpdated.Id))
	os.RemoveAll(configBackup.BACKUP_SRC)
	main.shutdowner.Shutdown()
}

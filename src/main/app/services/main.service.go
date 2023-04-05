package mainService

import (
	"fmt"
	configService "go-rotate-backup-s3/commons/app/services/config-service"
	logService "go-rotate-backup-s3/commons/app/services/log-service"
	filesDomain "go-rotate-backup-s3/commons/domain/files"
	filesService "go-rotate-backup-s3/files/app/services"
	"os"
	"os/exec"
	"time"
)

type MainService struct {
	fileService *filesService.FilesService
}

func New(fileService *filesService.FilesService) *MainService {
	return &MainService{fileService: fileService}
}

func (main *MainService) Run() {
	configBackup := configService.GetBackup()

	// command := strings.Split(configBackup.BACKUP_COMMAND, " ")

	// cmd := exec.Command(command[0], command[1:]...)
	cmd := exec.Command("bash", "-c", configBackup.BACKUP_COMMAND)
	out, err := cmd.CombinedOutput()

	if err != nil {
		logService.Error("executed command failed")
		logService.Error(err.Error())
		panic(1)
	}

	logService.Info("executed command successfull")
	fmt.Println(string(out))

	file, err := os.Open(configBackup.BACKUP_SRC)

	if err != nil {
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
			// fmt.Printf("%+v \n", v.Id)
			main.fileService.Delete(v.Id)
		}

		for _, v := range srcListFolder {
			// fmt.Printf("%+v \n", v.Id)
			main.fileService.Copy(v.Id, dst+v.Name)
			main.fileService.Delete(v.Id)
		}

		keep = keep - 1
	}

	fileUpdated, err := main.fileService.UploadLarge(*fileDomain)

	if err != nil {
		logService.Error(err.Error())
	}

	logService.Info(fmt.Sprintf("create successfull backup: %s", fileUpdated.Id))
}

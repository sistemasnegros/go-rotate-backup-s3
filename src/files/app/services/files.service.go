package filesService

import (
	filesDomain "go-rotate-backup-s3/commons/domain/files"
	s3InfraRepo "go-rotate-backup-s3/files/infra/s3"
)

type FilesService struct {
	repo *s3InfraRepo.S3Repository
}

func New(repo *s3InfraRepo.S3Repository) *FilesService {
	return &FilesService{repo: repo}
}

func (s *FilesService) Upload(file filesDomain.File) (*filesDomain.FileRes, error) {
	res, err := s.repo.Upload(file)
	return res, err

}

func (s *FilesService) Copy(src string, dst string) (*filesDomain.FileRes, error) {
	res, err := s.repo.Copy(src, dst)
	return res, err
}

func (s *FilesService) UploadLarge(file filesDomain.File) (*filesDomain.FileRes, error) {
	res, err := s.repo.UploadLarge(file)
	return res, err

}

func (s FilesService) List(prefix string) ([]*filesDomain.FileRes, error) {
	res, err := s.repo.List(prefix)
	return res, err
}

func (s FilesService) Get(id string) (*filesDomain.FileRes, error) {
	res, err := s.repo.Get(id)
	return res, err
}

func (s FilesService) Delete(id string) error {
	err := s.repo.Delete(id)
	return err
}

func (s FilesService) Download(id string) (*filesDomain.File, error) {
	res, err := s.repo.Download(id)
	return res, err
}

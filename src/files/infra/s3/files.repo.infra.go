package s3InfraRepo

import (
	"context"
	"fmt"
	configService "go-rotate-backup-s3/commons/app/services/config-service"
	logService "go-rotate-backup-s3/commons/app/services/log-service"
	filesDomain "go-rotate-backup-s3/commons/domain/files"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/exp/slices"
)

type S3Repository struct {
	Client *s3.Client
	Bucket string
	Url    string
}

func New(Client *s3.Client) *S3Repository {
	s3Config := configService.GetS3()
	return &S3Repository{
		Bucket: s3Config.AWS_BUCKET,
		Client: Client,
		Url:    s3Config.AWS_URL_PUBLIC,
	}
}

func (r *S3Repository) List(prefix string) ([]*filesDomain.FileRes, error) {

	res, err := r.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(r.Bucket),
		Prefix: aws.String(prefix),
	})

	// files := []*filesDomain.FileRes{}
	files := make([]*filesDomain.FileRes, len(res.Contents))

	if err != nil {
		logService.Error(err.Error())
		return files, err
	}

	for index, object := range res.Contents {
		key := aws.ToString(object.Key)
		paths := strings.Split(key, "/")
		name := paths[len(paths)-1]
		file := &filesDomain.FileRes{
			Id:   key,
			Name: name,
			Url:  fmt.Sprintf("%s/%s", r.Url, key),
		}

		files[index] = file

	}

	return files, nil
}

func (r *S3Repository) Upload(file filesDomain.File) (*filesDomain.FileRes, error) {

	_, err := r.Client.PutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(r.Bucket),
			Key:    aws.String(file.Id),
			Body:   file.Data,
			ACL:    "public-read",
		})

	fileUpdated := &filesDomain.FileRes{}

	if err != nil {
		logService.Error(err.Error())
		return fileUpdated, err
	}

	path := strings.Split(file.Id, "/")
	name := path[len(path)-1]

	fileUpdated.Id = file.Id
	fileUpdated.Name = name
	fileUpdated.Url = fmt.Sprintf("%s/%s", r.Url, file.Id)

	return fileUpdated, err

}

func (r *S3Repository) Copy(src string, dst string) (*filesDomain.FileRes, error) {
	_, err := r.Client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(r.Bucket),
		CopySource: aws.String(fmt.Sprintf("%v/%v", r.Bucket, src)),
		Key:        aws.String(fmt.Sprintf("%v", dst)),
	})

	if err != nil {
		logService.Error(fmt.Sprintf("Couldn't copy object from %v:%v to %v:%v. Here's why: %v\n",
			r.Bucket, src, r.Bucket, dst, err.Error()))
	}

	fileUpdated := &filesDomain.FileRes{}

	path := strings.Split(dst, "/")
	name := path[len(path)-1]

	fileUpdated.Id = dst
	fileUpdated.Name = name
	fileUpdated.Url = fmt.Sprintf("%s/%s", r.Url, dst)

	return fileUpdated, err

}

func (r *S3Repository) UploadLarge(file filesDomain.File) (*filesDomain.FileRes, error) {
	var partMiBs int64 = 10
	uploader := manager.NewUploader(r.Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(file.Id),
		Body:   file.Data,
		// ACL:    "public-read",
	})

	fileUpdated := &filesDomain.FileRes{}

	if err != nil {
		logService.Error(fmt.Sprintf("Couldn't upload large object to %v:%v. Here's why: %v\n",
			r.Bucket, file.Id, err))
	}

	path := strings.Split(file.Id, "/")
	name := path[len(path)-1]

	fileUpdated.Id = file.Id
	fileUpdated.Name = name
	fileUpdated.Url = fmt.Sprintf("%s/%s", r.Url, file.Id)

	return fileUpdated, err
}

func (r *S3Repository) Download(id string) (*filesDomain.File, error) {

	res, err := r.Client.GetObject(context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(r.Bucket),
			Key:    aws.String(id),
		})

	file := &filesDomain.File{}

	if err != nil {
		logService.Error(err.Error())
		return file, err
	}

	path := strings.Split(id, "/")
	name := path[len(path)-1]

	file.Id = id
	file.Name = name
	file.ContentType = *res.ContentType
	file.Data = res.Body
	defer res.Body.Close()

	return file, err
}

func (r *S3Repository) Get(id string) (*filesDomain.FileRes, error) {

	res, err := r.List("/")

	if err != nil {
		logService.Error(err.Error())
	}

	index := slices.IndexFunc(res, func(f *filesDomain.FileRes) bool { return f.Id == id })

	return res[index], err
}

func (r *S3Repository) Delete(id string) error {

	_, err := r.Client.DeleteObject(context.TODO(),
		&s3.DeleteObjectInput{
			Bucket: aws.String(r.Bucket),
			Key:    aws.String(id),
		})

	if err != nil {
		logService.Error(err.Error())
	}

	logService.Info(fmt.Sprintf("file deleted successfull %s", id))

	return err
}

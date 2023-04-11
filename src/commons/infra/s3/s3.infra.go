package s3Infra

import (
	"context"

	configService "go-rotate-backup-s3/commons/app/services/config-service"
	logService "go-rotate-backup-s3/commons/app/services/log-service"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func New() *s3.Client {
	s3Config := configService.GetS3()

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           s3Config.AWS_ENDPOINT,
			SigningRegion: s3Config.AWS_REGION,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3Config.AWS_ACCESS_KEY_ID, s3Config.AWS_SECRET_ACCESS_KEY, "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		logService.Error(err.Error())
		panic(err)
	}

	logService.Info("connection S3 successful")
	client := s3.NewFromConfig(cfg)

	return client
}

package tilefetcher

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type MyProvider struct{}

func (m *MyProvider) Retrieve() (credentials.Value, error) {
	AccessKeyID := os.Getenv("KEY")
	SecretAccessKey := os.Getenv("SECRET")

	value := credentials.Value{
		AccessKeyID:     AccessKeyID,
		SecretAccessKey: SecretAccessKey,
	}
	return value, nil
}

func (m *MyProvider) IsExpired() bool {
	return false
}

func Fetch(path string) (string, error) {
	creds := credentials.NewCredentials(&MyProvider{})

	// Specify the region
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(endpoints.EuWest1RegionID),
	}))

	downloader := s3manager.NewDownloader(sess)

	filepath := path + "downloaded-tile"
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %q, %v", file, err)
	}

	fmt.Println("Downloading tile. This could take a while...")
	n, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String("pcf-rabbitmq-pipelines"),
		Key:    aws.String("rabbitmq-1.18/archive/p-rabbitmq-1.18.1-build.30.pivotal"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to download file, %v, err")
	}

	fmt.Printf("file downloaded, %d bytes\n", n)

	return filepath, nil
}

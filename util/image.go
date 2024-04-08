package util

import (
	"fmt"
	_ "image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var baseURL string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	baseURL = fmt.Sprintf("https://%s.s3.amazonaws.com/", os.Getenv("BUCKET_NAME"))
}

func SaveImage(src io.ReadSeeker) (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s.jpg", uid.String())

	credential := credentials.NewStaticCredentials(
		os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_KEY"),
		"",
	)

	awsConfig := aws.Config{
		Region:      aws.String(os.Getenv("REGION")),
		Credentials: credential,
	}

	s, err := session.NewSession(&awsConfig)
	if err != nil {
		return "", err
	}

	svc := s3.New(s)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(fileName),
		Body:   src,
	})
	if err != nil {
		return "", err
	}

	return filepath.Join(baseURL, fileName), nil
}

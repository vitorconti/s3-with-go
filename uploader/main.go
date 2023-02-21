package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	config "github.com/vitorconti/s3-with-go/config"
)

var (
	s3Client *s3.S3
	s3Bucket string
	wg       sync.WaitGroup
)

func init() {
	configs, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	sess, err := session.NewSession(
		&aws.Config{
			Region: &configs.AwsRegion,
			Credentials: credentials.NewStaticCredentials(
				configs.AwsKey, configs.AwsPassword, ""),
		},
	)
	s3Client = s3.New(sess)
	s3Bucket = configs.AwsBucketName
}
func main() {
	dir, err := os.Open("./tmp")
	if err != nil {
		panic(err)
	}
	defer dir.Close()
	uploadControl := make(chan struct{}, 100)
	errorFileUpload := make(chan string, 10)

	go func() {
		for {
			select {
			case filename := <-errorFileUpload:
				uploadControl <- struct{}{}
				wg.Add(1)
				go uploadFile(filename, uploadControl, errorFileUpload)
			}
		}
	}()

	for {
		files, err := dir.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading directory: %s\n", err)
			continue
		}
		wg.Add(1)
		uploadControl <- struct{}{}
		go uploadFile(files[0].Name(), uploadControl, errorFileUpload)
	}
	wg.Wait()

}

func uploadFile(filename string, uploadControl <-chan struct{}, errorFileUpload <-chan string) {
	completeFileName := fmt.Sprintf("./tmp/%s", filename)
	fmt.Printf("Uploading file %s to bucket %s \n", completeFileName, s3Bucket)
	f, err := os.Open(completeFileName)
	if err != nil {
		fmt.Printf("Error opening file %s \n", completeFileName)
		<-uploadControl
		<-errorFileUpload
		return
	}
	defer f.Close()
	_, err = s3Client.PutObject(
		&s3.PutObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(filename),
			Body:   f,
		},
	)
	if err != nil {
		fmt.Printf("Error uploading file %s \n", completeFileName)
		<-uploadControl
		<-errorFileUpload
		return
	}
	fmt.Printf("File uploaded successfully %s \n", completeFileName)
	<-uploadControl
}

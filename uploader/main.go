package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vitorcontiz/s3-with-go/config"
)

var (
	s3Client *s3.S3
	s3Bucket string
)

func init() {
	configs,err := config.LoadConfig()
	sess, err := session.NewSession(
		&aws.Config{
			Region: configs.Region,
		}
	)
}
func main() {

}

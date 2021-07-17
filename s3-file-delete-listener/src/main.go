package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

var (
	region     = os.Getenv("aws_region")
	bucketName = os.Getenv("s3_bucket_name")
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Println("creation of session error: ", err.Error())
	}
	svc := s3.New(sess)
	for _, record := range snsEvent.Records {
		log.Println("File name to be deleted: " + record.SNS.Message)
		resp, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(record.SNS.Message),
		})
		if err != nil {
			log.Println("Panic Error in deleting file for S3: ", err.Error())
		}
		log.Println("RESP", resp)
	}
}

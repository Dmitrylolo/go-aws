package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	fmt.Println("Creating buckets")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	svc := s3.New(sess)

	input := &s3.CreateBucketInput{
		Bucket: aws.String("go-s3-example-bucket"),
	}

	result, err := svc.CreateBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				log.Fatal("Bucket is owned by you")
				break
			case s3.ErrCodeBucketAlreadyExists:
				log.Fatal("Bucket already exists")
				break
			default:
				log.Fatal(aerr.Error())
			}
		} else {
			log.Fatal(err.Error())
		}
	}

	log.Print(result)
	// buckets, err := svc.ListBuckets(&s3.ListBucketsInput{})
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
}

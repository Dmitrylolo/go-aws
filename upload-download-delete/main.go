package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadItem(sess *session.Session) {
	f, err := os.Open("testtext.txt")
	if err != nil {
		log.Fatal("could not open file")
	}
	defer f.Close()

	uploader := s3manager.NewUploader(sess)
	result, error := uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String("public-read"),
		Bucket: aws.String("go-s3-example-bucket"),
		Key:    aws.String("test.txt"),
		Body:   f,
	})
	if error != nil {
		log.Fatal(error.Error())
	}

	log.Printf("Upload Result: %v\n", result)
}

func listItems(sess *session.Session) {
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("go-s3-example-bucket"),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, item := range resp.Contents {
		log.Printf("Name: %s\n", *item.Key)
		log.Printf("LastModified: %s\n", *item.LastModified)
		log.Printf("Size: %d\n", *item.Size)
	}
}

func downloadItem(sess *session.Session) {
	file, err := os.Create("downloaded.txt")
	if err != nil {
		log.Fatal("could not create file")
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String("go-s3-example-bucket"),
			Key:    aws.String("test.txt"),
		})
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Download complete")
}

func deleteItem(sess *session.Session) {
	svc := s3.New(sess)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String("go-s3-example-bucket"),
		Key:    aws.String("test.txt"),
	}

	result, err := svc.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				log.Println("The specified key does not exist.")
			default:
				log.Fatal(err.Error())
			}
		} else {
			log.Fatal(err.Error())
		}
	}

	log.Printf("Successfully deleted object %+v\n", result)
}

func main() {
	fmt.Println("Upload and download")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	UploadItem(sess)
	listItems(sess)
	downloadItem(sess)
	deleteItem(sess)
	listItems(sess)
}

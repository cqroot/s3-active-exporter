package internal

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"

	"github.com/cqroot/s3-active-exporter/logger"
)

var (
	bucket    string
	accessKey string
	secretKey string
	once      sync.Once
)

type S3ActiveMonitor struct {
	sess     *session.Session
	endpoint string
	filename string
}

func (m *S3ActiveMonitor) Run(endpoint string) (
	putResult float64, getResult float64, delResult float64,
	putFilename string, getFilename string, delFilename string) {
	once.Do(func() {
		bucket = viper.GetString("s3.bucket")
		accessKey = viper.GetString("s3.access")
		secretKey = viper.GetString("s3.secret")
	})
	m.endpoint = endpoint
	m.filename = fmt.Sprintf("%s-%s", m.endpoint, time.Now().Format("2006_01_02_15_04_05"))

	putResult = 0
	getResult = 0
	delResult = 0
	putFilename = ""
	getFilename = ""
	delFilename = ""

	if bucket == "" || accessKey == "" || secretKey == "" {
		logger.Error("you must supply a bucket name, access key and secret key")
		putFilename = "you must supply a bucket name, access key and secret key"
		getFilename = "you must supply a bucket name, access key and secret key"
		delFilename = "you must supply a bucket name, access key and secret key"
		return
	}

	var err error
	m.sess, err = session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
		Endpoint:         aws.String(m.endpoint),
		Region:           aws.String(viper.GetString("s3.region")),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(true),
	})
	if err != nil {
		logger.Error(err.Error())
		putFilename = err.Error()
		getFilename = err.Error()
		delFilename = err.Error()
	}

	putFilename, putResult = m.putObject()
	getFilename, getResult = m.getObject()
	delFilename, delResult = m.delObject()
	return
}

func (m *S3ActiveMonitor) putObject() (string, float64) {
	uploader := s3manager.NewUploader(m.sess)

	data := "Hello, world!"
	reader := strings.NewReader(data)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(m.filename),
		Body:   reader,
	})
	if err != nil {
		logger.Error(err.Error())
		return err.Error(), 0
	}

	return m.filename, 1
}

func (m *S3ActiveMonitor) getObject() (string, float64) {
	downloader := s3manager.NewDownloader(m.sess)

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(m.filename),
		})
	if err != nil {
		logger.Error(err.Error())
		return err.Error(), 0
	}

	return m.filename, 1
}

func (m *S3ActiveMonitor) delObject() (string, float64) {
	svc := s3.New(m.sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(m.filename)})
	if err != nil {
		logger.Error(err.Error())
		return err.Error(), 0
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(m.filename),
	})
	if err != nil {
		logger.Error(err.Error())
		return err.Error(), 0
	}

	return m.filename, 1
}
